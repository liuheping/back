package context

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppName    string
	ListenPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret   string
	JWTExpireIn time.Duration

	DebugMode bool
	LogFormat string

	QiniuBucket    string
	QiniuAccessKey string
	QiniuSecretKey string
	QiniuCndHost   string

	AdminADPositionOptions    string
	MerchantADPositionOptions string
}

func LoadConfig(path string) *Config {
	config := viper.New()
	config.SetConfigName("Config")
	config.AddConfigPath(".")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error context file: %s \n", err)
	}

	return &Config{
		AppName:    config.Get("app-name").(string),
		ListenPort: config.Get("listen-port").(string),

		DBHost:     config.Get("db.host").(string),
		DBPort:     config.Get("db.port").(string),
		DBUser:     config.Get("db.user").(string),
		DBPassword: config.Get("db.password").(string),
		DBName:     config.Get("db.dbname").(string),

		JWTSecret:   config.Get("auth.jwt-secret").(string),
		JWTExpireIn: config.GetDuration("auth.jwt-expire-in"),

		DebugMode: config.Get("log.debug-mode").(bool),
		LogFormat: config.Get("log.log-format").(string),

		QiniuAccessKey: config.GetString("qiniu.AccessKey"),
		QiniuSecretKey: config.GetString("qiniu.SecretKey"),
		QiniuBucket:    config.GetString("qiniu.bucket"),
		QiniuCndHost:   config.GetString("qiniu.cndHost"),

		AdminADPositionOptions:    config.GetString("ad.admin_position_options"),
		MerchantADPositionOptions: config.GetString("ad.merchant_position_options"),
	}
}
