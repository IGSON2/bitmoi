package utilities

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	SymmetricKey string `mapstructure:"SYMMETRIC_KEY"`
}

var C *Config

func GetConfig() *Config {
	if C == nil {
		viper.AddConfigPath("../../")
		viper.SetConfigName("app")
		viper.SetConfigType("env")

		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			log.Panicln("Err! cannot load config. Err : ", err)
		}

		if err := viper.Unmarshal(&C); err != nil {
			log.Panicln("Err! cannot unmarshal config. Err : ", err)
		}
	}
	return C
}
