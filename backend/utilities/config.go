package utilities

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	SymmetricKey         string        `mapstructure:"SYMMETRIC_KEY"`
	DBDriver             string        `mapstructure:"DB_Driver"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	HTTPAddress          string        `mapstructure:"HTTP_ADDRESS"`
	GRPCAddress          string        `mapstructure:"GRPC_ADDRESS"`
	apikey               string        `mapstructure:"API_KEY"`
	srckey               string        `mapstructure:"SRC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

var C *Config

func GetConfig(path string) *Config {
	if C == nil {
		viper.AddConfigPath(path)
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
	cfg := GetConfig("../../.")
	return cfg.apikey, cfg.srckey
}
