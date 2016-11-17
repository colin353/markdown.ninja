package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/imdario/mergo"

	"gopkg.in/yaml.v2"
)

// EnvPrefix is the prefix appended to the configuration field names
// in their environment variable override names. So, for example, the
// YAML RedisURL is overridden by the environment variable
// APPCONFIG_REDISURL.
var EnvPrefix = "APPCONFIG_"

// Config is a struct which contains all of the configuration for
// the app.
type Config struct {
	Port          string
	RedisURL      string
	DataDirectory string
	Hostnames     []string
}

// LoadConfig generates the configuration using three rules:
// 1. get the default configuration from default.yaml
// 2. get the special configuration from config.yaml
// 3. override all configuration parameters with enviroment vars
//    if they exist.
func LoadConfig(directory string) *Config {
	config := Config{}
	defaultConfig := Config{}

	defaultConfigBytes, err := ioutil.ReadFile(directory + "/default.yaml")
	if err != nil {
		panic("Unable to read configuration file.")
	}
	err = yaml.Unmarshal(defaultConfigBytes, &defaultConfig)
	if err != nil {
		panic("Unable to parse configuration file `default.yaml`.")
	}

	// Try to load special config from the config.yaml file, but if
	// it doesn't exist, that's fine, we can ignore it.
	specialConfigBytes, err := ioutil.ReadFile(directory + "/config.yaml")
	if err == nil {
		err = yaml.Unmarshal(specialConfigBytes, &config)
		if err != nil {
			panic("Unable to parse configuration file `config/default.yaml`.")
		}
	}

	mergo.Merge(&config, defaultConfig)

	// Now override all parameters with environment variables, if they exist.
	instanceValue := reflect.ValueOf(&config).Elem()
	instanceType := reflect.TypeOf(&config).Elem()
	for i := 0; i < instanceValue.NumField(); i++ {
		t := instanceType.Field(i)
		v := instanceValue.Field(i)

		env := os.Getenv(EnvPrefix + strings.ToUpper(t.Name))
		if env == "" {
			continue
		}
		// If we reach this point, the user has specified an environment variable
		// override. So we need to check the destination type, and override the
		// config value with the environment variable.

		typ := v.Type().Kind()
		switch typ {
		case reflect.String:
			v.SetString(env)
		// Only an array/slice of strings is permitted right now.
		case reflect.Slice:
			array := strings.Split(env, ",")
			value := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf("")), 0, len(array))

			for _, el := range array {
				x := reflect.ValueOf(el)
				value = reflect.Append(value, x)
			}
			v.Set(value)
		default:
			log.Printf("Unexpected type in configuration: %v", typ)
		}
	}

	result, err := json.MarshalIndent(config, "", "  ")
	log.Printf("Starting up with settings: \n %s", string(result))

	return &config
}
