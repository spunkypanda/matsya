package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

var config map[string]string

// GetEnvironmentName receives a string and returns an environment name
func GetEnvironmentName(env string) string {
	env = os.Getenv(env)
	if env == "" {
		env = "development"
	}
	return env
}

// Initialize initializes config management
func Initialize(ENV, prefixPath string) {
	configName := ENV
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(prefixPath + "src/config/")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}

// GetConfigMap Fetch config map
func GetConfigMap(key string) map[string]string {
	return viper.GetStringMapString(key)
}

// GetStringList Fetch config list for a key
func GetStringList(key string) []string {
	return viper.GetStringSlice(key)
}

// GetString Fetch config value for key
func GetString(key string) string {
	return viper.GetString(key)
}

// GetFlag Fetch config flag
func GetFlag(key string) bool {
	return viper.GetBool(key)
}
