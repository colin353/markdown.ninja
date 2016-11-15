package config

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

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
}

// LoadConfig generates the configuration using three rules:
// 1. get the default configuration from default.yaml
// 2. get the special configuration from config.yaml
// 3. override all configuration parameters with enviroment vars
//    if they exist.
func LoadConfig(directory string) *Config {
	c := Config{}

	defaultConfig, err := ioutil.ReadFile(directory + "/default.yaml")
	if err != nil {
		panic("Unable to read configuration file.")
	}
	err = yaml.Unmarshal(defaultConfig, &c)
	if err != nil {
		panic("Unable to parse configuration file `default.yaml`.")
	}

	// Try to load special config from the config.yaml file, but if
	// it doesn't exist, that's fine, we can ignore it.
	specialConfig, err := ioutil.ReadFile(directory + "/config.yaml")
	if err == nil {
		err = yaml.Unmarshal(specialConfig, &c)
		if err != nil {
			panic("Unable to parse configuration file `config/default.yaml`.")
		}
	}

	// Now override all parameters with environment variables, if they exist.
	instanceValue := reflect.ValueOf(&c).Elem()
	instanceType := reflect.TypeOf(&c).Elem()
	for i := 0; i < instanceValue.NumField(); i++ {
		t := instanceType.Field(i)
		v := instanceValue.Field(i)

		env := os.Getenv(EnvPrefix + strings.ToUpper(t.Name))
		if env != "" {
			log.Printf("Override detected on %v", EnvPrefix+strings.ToUpper(t.Name))
			v.SetString(env)
		}
	}
	return &c
}
