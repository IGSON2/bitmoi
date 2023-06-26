package utilities

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	SymmetricKey         string        `mapstructure:"SYMMETRIC_KEY"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	HTTPAddress          string        `mapstructure:"HTTP_ADDRESS"`
	GRPCHTTPAddress      string        `mapstructure:"GRPC_GATEWAY_ADDRESS"`
	GRPCAddress          string        `mapstructure:"GRPC_ADDRESS"`
	apikey               string        `mapstructure:"API_KEY"`
	srckey               string        `mapstructure:"SRC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	DataDir              string
}

var C *Config

func (c *Config) SetDataDir(path string) {
	c.DataDir = path
}

func GetConfig(path string) *Config {
	if C == nil {
		viper.AddConfigPath(path)
		viper.SetConfigName("a")
		viper.SetConfigType("env")

		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			matches, _ := filepath.Glob(path)
			files := strings.Join(matches, ",")
			log.Panicf("Err! cannot load config. Err : %v, File list : %s", err, files)
		}

		if err := viper.Unmarshal(&C); err != nil {
			log.Panicln("Err! cannot unmarshal config. Err : ", err)
		}
		C.SetDataDir(DefaultDataDir())
	}
	return C
}

func GetAPIKeys() (apikey, srckey string) {
	cfg := GetConfig("../../.")
	return cfg.apikey, cfg.srckey
}
