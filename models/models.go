package models

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/colin353/markdown.ninja/config"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

// AppConfig is an instance of the application config.
var AppConfig *config.Config

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
	RegistrationKey() string
	Export() map[string]interface{}
}

// A ModelIterator is an iterator over a list of models.
type ModelIterator interface {
	Next() bool
	Value() Model
	Count() int
}

var connectionPool *pool.Pool

// This function gets and returns a redis client object. If we are
// currently in test mode, we'll make sure to select the test database,
// which is redis database 1. Otherwise, we're in production, so we'll
// use the permanent database 0 (default).
func getRedisConnection() (*redis.Client, error) {
	p, err := connectionPool.Get()
	if err != nil {
		return nil, err
	}
	if AppConfig.Mode == "test" || AppConfig.Mode == "testing" {
		response := p.Cmd("SELECT", 1)
		if response.Err != nil {
			log.Fatalf("Unable to select testing database: %v", response.Err.Error())
			return nil, response.Err
		}
	}

	return p, err
}

// Connect to the redis database and create a connection pool, which can be used
// to query the redis database concurrently.
func Connect() {
	var err error
	connectionPool, err = pool.New("tcp", AppConfig.RedisURL, 10)
	if err != nil {
		log.Fatal("Unable to connect to redis server.")
	}
}

// ClearDatabase deletes all the keys in the database. As a precautionary
// measure, it also selects database 1, which is designated as the testing
// database.
func ClearDatabase() {
	p, err := getRedisConnection()
	if err != nil {
		log.Fatal("Couldn't connect to the redis database.")
	}
	response := p.Cmd("SELECT", 1)
	if response.Err != nil {
		log.Fatal("Unable to select testing database, terminating.")
	}

	// Clear all keys in the database. This command cannot fail,
	// according to redis docs.
	p.Cmd("FLUSHDB")
}

// MakeKeyForTable creates a unique key for a table by using the redis INCR
// command on the key:tablename field. If you aren't sure what key to give
// a new object, use MakeKeyForTable and give a string for the table name,
// and you'll be guaranteed a unique key.
func MakeKeyForTable(table string) string {
	p, err := getRedisConnection()
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
	p, err := getRedisConnection()
	if err != nil {
		log.Fatal("Couldn't connect to the redis database.")
		return err
	}

	// Delete the hash.
	result := p.Cmd("DEL", m.Key())
	if result.Err != nil {
		return result.Err
	}

	// Also, delete the member from the set.
	result = p.Cmd("SREM", m.RegistrationKey(), m.Key())
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
		return errors.New("proposed model changes failed to validate")
	}

	// Okay, the changes are fine to apply to the database.
	p, err := getRedisConnection()
	if err != nil {
		log.Fatal("couldn't connect to the redis database")
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
	err := saveOrInsert(m, false)
	if err != nil {
		return err
	}

	// Now, register the RegistrationKey into the registration set.
	p, err := getRedisConnection()
	if err != nil {
		log.Fatal("Couldn't connect to the redis database.")
		return err
	}
	p.Cmd("SADD", m.RegistrationKey(), m.Key())

	return nil
}

// ModelList is an implementation of a ModelIterator.
type ModelList struct {
	Prototype Model
	Keys      []string
	Index     int
}

// Next loads the next value in the iterator, and returns
// true if it succeeded. To get the value, do m.Value().
func (m *ModelList) Next() bool {
	if m.Index >= len(m.Keys) {
		return false
	}
	err := LoadFromKey(m.Prototype, m.Keys[m.Index])
	m.Index++

	// If there is an error loading a key, skip ahead.
	if err != nil {
		return m.Next()
	}
	return true
}

// Value returns the currently-pointed-to value for the
// iterator.
func (m *ModelList) Value() Model {
	return m.Prototype
}

// Count returns the total number of objects available
// in the iterator.
func (m *ModelList) Count() int {
	return len(m.Keys)
}

// GetList takes an object and returns a list of sibilings of it.
// Actually it returns a ModelIterator, which you can call "Next()" on
// any number of times.
func GetList(m Model) (ModelIterator, error) {
	p, err := getRedisConnection()
	if err != nil {
		log.Fatal("couldn't connect to the redis database")
		return nil, err
	}

	result := p.Cmd("SMEMBERS", m.RegistrationKey())

	keys, err := result.List()
	if err != nil {
		return nil, err
	}

	return &ModelList{Prototype: m, Keys: keys}, result.Err
}

func saveOrInsert(m Model, expectKey bool) error {
	// Ensure that the model is validated.
	if !m.Validate() {
		return errors.New("model failed to validate")
	}

	p, err := getRedisConnection()
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
		return fmt.Errorf("Key `%s` already exists: can't insert. Did you mean to save?", m.Key())
	}
	if expectKey && keyAlreadyExists == 0 {
		return fmt.Errorf("Key `%s` doesn't exist: can't save. Did you mean to insert?", m.Key())
	}

	response = p.Cmd("HMSET", m.Key(), instanceMap)
	if response.Err != nil {
		log.Fatal("Error executing redis save command.")
		return err
	}

	return nil
}

// Load takes a partially filled out Model struct, searches for it in the
// database, and fills in the rest of the fields.
func Load(m Model) error {
	return LoadFromKey(m, m.Key())
}

// LoadFromKey takes a model and a key string,
// and searches in the redis database for that object, then fills out
// the object fields with whatever is in the database. For any fields
// which don't exist in the database, it uses m.MakeDefault() to set
// fields to their default values.
func LoadFromKey(m Model, key string) error {
	m.MakeDefault()

	p, err := getRedisConnection()
	if err != nil {
		log.Print("couldn't connect to the redis database")
		return err
	}
	response := p.Cmd("HGETALL", key)
	if err != nil {
		log.Print("error executing redis load command (HGETALL)")
		return err
	}

	instanceMap, err := response.Map()
	if err != nil {
		log.Printf("Unable to unwind key '%s' in redis as a map (!)", key)
		return err
	}
	// Check if there are no keys in the map: that means it doesn't exist.
	if len(instanceMap) == 0 {
		return fmt.Errorf("there is no such key: `%v`", key)
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
		if fieldName == "" {
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
