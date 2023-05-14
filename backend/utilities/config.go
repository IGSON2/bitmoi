package utilities

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	SymmetricKey string `mapstructure:"SYMMETRIC_KEY"`
	DBDriver     string `mapstructure:"DB_Driver"`
	DBSource     string `mapstructure:"DB_SOURCE"`
	Address      string `mapstructure:"ADDRESS"`
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
