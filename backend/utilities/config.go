package utilities

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	symmetricKey string `mapstructure:"SYMMETRIC_KEY"`
}

var config *Config

func init() {
	viper.AddConfigPath("../../")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Panicln("Err! cannot load config. Err : ", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Panicln("Err! cannot unmarshal config. Err : ", err)
	}
}

func GetConfig() *Config {
	return config
}
