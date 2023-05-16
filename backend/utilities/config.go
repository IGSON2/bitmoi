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
	apikey       string `mapstructure:"apikey"`
	srckey       string `mapstructure:"srckey"`
}

var C *Config

func GetConfig() *Config {
	if C == nil {
		viper.AddConfigPath("../../")
		viper.SetConfigName("a")
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

func GetAPIKeys() (apikey, srckey string) {
	cfg := GetConfig()
	return cfg.apikey, cfg.srckey
}
