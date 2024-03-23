package config

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"bitmoi/backend/utilities/common"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Config is the struct for the configuration of the application
type Config struct {
	Environment             string        `mapstructure:"ENVIRONMENT"`
	SymmetricKey            string        `mapstructure:"SYMMETRIC_KEY"`
	DBDriver                string        `mapstructure:"DB_DRIVER"`
	DBSource                string        `mapstructure:"DB_SOURCE"`
	HTTPAddress             string        `mapstructure:"HTTP_ADDRESS"`
	GRPCHTTPAddress         string        `mapstructure:"GRPC_GATEWAY_ADDRESS"`
	GRPCAddress             string        `mapstructure:"GRPC_ADDRESS"`
	RedisAddress            string        `mapstructure:"REDIS_ADDRESS"`
	AccessTokenDuration     time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration    time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	EmailSenderPassword     string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	PrivateKey              string        `mapstructure:"PRIVATE_KEY"`
	S3AccessKey             string        `mapstructure:"S3_ACCESS_KEY"`
	S3SecretKey             string        `mapstructure:"S3_SECRET_KEY"`
	BiddingDuration         time.Duration `mapstructure:"BIDDING_DURATION"`
	GoogleOauthClientID     string        `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleOauthClientSecret string        `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
	KakaoOauthClientID      string        `mapstructure:"KAKAO_OAUTH_CLIENT_ID"`
	OauthRedirectURL        string
	DataDir                 string
	LogLevel                zerolog.Level
}

var C *Config

func (c *Config) SetDataDir(path string) {
	c.DataDir = path
}

func (c *Config) SetLogLevel(level int8) {
	c.LogLevel = zerolog.Level(level)
}

func (c *Config) SwitchOauthRedirectURL() {
	if c.Environment == common.EnvProduction {
		c.OauthRedirectURL = "https://m.bitmoi.co.kr"
	} else {
		c.OauthRedirectURL = "http://localhost:3000"
	}
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
		C.SwitchOauthRedirectURL()
	}
	return C
}
