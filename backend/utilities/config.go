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
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	EmailSenderPassword  string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	PrivateKey           string        `mapstructure:"PRIVATE_KEY"`
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
		matches, _ := filepath.Glob(path)
		files := strings.Join(matches, ",")

		if err := viper.ReadInConfig(); err != nil {
			log.Panicf("Err! cannot load config. Err : %v, File list : %s", err, files)
		}

		if err := viper.Unmarshal(&C); err != nil {
			log.Panicf("Err! cannot unmarshal config. Err : %v, File list : %s", err, files)
		}
		C.SetDataDir(DefaultDataDir())
	}
	return C
}
