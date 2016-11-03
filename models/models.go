package models

import (
	"errors"
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	"log"
	"reflect"
	"strconv"
)

// Model is the basic interface required for an object to be
// saved to the redis store. Fundamentally, you need a function
// which generates a unique key for that individual structure
// (typically prefixed by the table name, for example tablename:identifier)
// and a validation function to check if the data is correct. Also, you need
// a function to set all fields to their defaults in case the fields don't
// exist in the database (e.g. due to migration).
type Model interface {
	Key() string
	MakeDefault()
	Validate() bool
}

var connectionPool *pool.Pool

// Connect to the redis database and create a connection pool, which can be used
// to query the redis database concurrently.
func init() {
	var err error
	connectionPool, err = pool.New("tcp", "localhost:6379", 10)
	if err != nil {
		log.Fatal("Unable to connect to redis server.")
	}
}

// MakeKeyForTable creates a unique key for a table by using the redis INCR
// command on the key:tablename field. If you aren't sure what key to give
// a new object, use MakeKeyForTable and give a string for the table name,
// and you'll be guaranteed a unique key.
func MakeKeyForTable(table string) string {
	p, err := connectionPool.Get()
	if err != nil {
		log.Fatal("Couldn't connect to the redis database.")
	}
	response := p.Cmd("INCR", fmt.Sprintf("key:%s", table))
	return response.String()
}

// This function takes a reflection value and sets the value based upon
// its internal type and the associated redis string.
func setFieldWithRedisString(v *reflect.Value, value string) {
	switch fieldType := v.Interface().(type) {
	case bool:
		result, err := strconv.ParseBool(value)
		if err != nil {
			log.Printf("Unable to convert %T while loading from redis model.\n", fieldType)
			return
		}
		v.SetBool(result)
	case int:
		result, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			log.Printf("Unable to convert %T while loading from redis model.\n", fieldType)
			return
		}
		v.SetInt(result)
	case string:
		v.SetString(value)
	default:
		panic(fmt.Sprintf("Unexpected type (%T) received while decoding redis model.", fieldType))
	}
}

// Delete removes a model from the database.
func Delete(m Model) error {
	p, err := connectionPool.Get()
	if err != nil {
		log.Fatal("Couldn't connect to the redis database.")
		return err
	}

	result := p.Cmd("DEL", m.Key())
	return result.Err
}

// UpdateWithChanges for when you only have a limited set of changes to make to the database,
// and don't want to do a read, apply updates, write process. For example, if you have a list
// of changes in a JSON object, which  may not encompass all fields in the object. This will
// check if those proposed changes are valid, then write them in a single shot.
func UpdateWithChanges(m Model, changes map[string]interface{}) error {
	// Start by updating the model struct with the changes, and validating that
	// the changes are legitimate by running m.Validate()
	instanceValue := reflect.ValueOf(m).Elem()
	instanceType := reflect.TypeOf(m).Elem()
	acceptedChanges := make(map[string]string)
	for i := 0; i < instanceValue.NumField(); i++ {
		t := instanceType.Field(i)
		v := instanceValue.Field(i)
		fieldName := t.Tag.Get("json")

		// Check if the field is in the list of proposed changes.
		updatedValue, ok := changes[fieldName]
		if !ok {
			continue
		}
		acceptedChanges[fieldName] = fmt.Sprintf("%v", updatedValue)
		setFieldWithRedisString(&v, acceptedChanges[fieldName])
	}

	if !m.Validate() {
		return errors.New("Proposed model changes failed to validate.")
	}

	// Okay, the changes are fine to apply to the database.
	p, err := connectionPool.Get()
	if err != nil {
		log.Fatal("Couldn't connect to the redis database.")
		return err
	}

	p.Cmd("HMSET", m.Key(), acceptedChanges)

	return nil
}

// Save takes an instance of a model and saves it to the database.
// If that model doesn't exist, it will return an error.
func Save(m Model) error {
	return saveOrInsert(m, true)
}

// Insert creates a new instance of a model in the database. It'll
// return an error if the key already exists.
func Insert(m Model) error {
	return saveOrInsert(m, false)
}

func saveOrInsert(m Model, expectKey bool) error {
	// Ensure that the model is validated.
	if !m.Validate() {
		return errors.New("Model failed to validate.")
	}

	p, err := connectionPool.Get()
	if err != nil {
		log.Fatal("Couldn't connect to the redis database.")
		return err
	}

	// Create a map[string]string reperesnting all the data in
	// the instance.
	instanceMap := make(map[string]string)
	instanceValue := reflect.ValueOf(m).Elem()
	instanceType := reflect.TypeOf(m).Elem()
	for i := 0; i < instanceValue.NumField(); i++ {
		t := instanceType.Field(i)
		v := instanceValue.Field(i)
		fieldName := t.Tag.Get("json")

		if fieldName == "" {
			continue
		}
		instanceMap[fieldName] = fmt.Sprintf("%v", v.Interface())
	}

	response := p.Cmd("EXISTS", m.Key())
	if response.Err != nil {
		log.Fatal("Error executing redis existence check command.")
		return err
	}

	keyAlreadyExists, _ := response.Int()
	if !expectKey && keyAlreadyExists == 1 {
		return errors.New("Key already exists: can't insert. Did you mean to save?")
	}
	if expectKey && keyAlreadyExists == 0 {
		return errors.New("Key doesn't exist: can't save. Did you mean to insert?")
	}

	response = p.Cmd("HMSET", m.Key(), instanceMap)
	if response.Err != nil {
		log.Fatal("Error executing redis save command.")
		return err
	}

	return nil
}

// Load takes a model (which must have a working Key() function)
// and searches in the redis database for that object, then fills out
// the object fields with whatever is in the database. For any fields
// which don't exist in the database, it uses m.MakeDefault() to set
// fields to their default values.
func Load(m Model) error {
	m.MakeDefault()

	p, err := connectionPool.Get()
	if err != nil {
		log.Print("Couldn't connect to the redis database.")
		return err
	}
	response := p.Cmd("HGETALL", m.Key())
	if err != nil {
		log.Print("Error executing redis load command (HGETALL).")
		return err
	}

	instanceMap, err := response.Map()
	if err != nil {
		log.Printf("Unable to unwind key '%s' in redis as a map (!)", m.Key())
		return err
	}
	// Check if there are no keys in the map: that means it doesn't exist.
	if len(instanceMap) == 0 {
		return errors.New("There is no such key.")
	}

	instanceValue := reflect.ValueOf(m).Elem()
	instanceType := reflect.TypeOf(m).Elem()
	// Iterate over the fields in the struct.
	for i := 0; i < instanceValue.NumField(); i++ {
		t := instanceType.Field(i)
		v := instanceValue.Field(i)
		fieldName := t.Tag.Get("json")

		// Only fields with JSON names will be loaded. Also, we use
		// the key field in the lookup, so it doesn't need to be loaded.
		if fieldName == "" || fieldName == "key" {
			continue
		}

		// Double-check that our instance map actually has a value for
		// this field. If not, then we'll skip the field.
		value, ok := instanceMap[fieldName]
		if !ok {
			continue
		}

		setFieldWithRedisString(&v, value)
	}

	return nil
}
