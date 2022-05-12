package config

import (
	"os"

	_ "gopkg.in/validator.v2" // validator
)

// GConf is Global config
var GConf = GetConfig()

// S3 is s3 config
type S3 struct {
	AccessKey string `validate:"nonnil"`
	SecureKey string `validate:"nonnil"`
	Bucket    string `validate:"nonnil"`
	Region    string `validate:"nonnil"`
	Url       string
}

// Config is configuraton for service
type Config struct {
	PostgresDatabase DBConfig
	S3               S3
}

type DBConfig struct {
	User     string `validate:"nonnil"`
	Password string `validate:"nonnil"`
	Database string `validate:"nonnil"`
	Host     string `validate:"nonnil"`
	Port     string `validate:"nonnil"`
	SSL      string `validate:"nonnil"`

	env string
}

// GetConfig gets config
func GetConfig() Config {

	// ------ postgres database -----
	var postgresDbConfig DBConfig
	postgresDbConfig.User = os.Getenv("DB_USER")
	postgresDbConfig.Database = os.Getenv("DB_DATABASE")
	postgresDbConfig.Host = os.Getenv("DB_HOST")
	postgresDbConfig.Password = os.Getenv("DB_PWD")
	postgresDbConfig.Port = os.Getenv("DB_PORT")
	postgresDbConfig.SSL = os.Getenv("DB_SSL")

	// ------ s3 -----
	s3Config := S3{}
	s3Config.AccessKey = os.Getenv("S3_ACCESS_KEY")
	s3Config.SecureKey = os.Getenv("S3_SECRET")
	s3Config.Bucket = os.Getenv("S3_BUCKET")
	s3Config.Region = os.Getenv("S3_REGION")
	s3Config.Url = os.Getenv("S3_URL")

	config := Config{
		PostgresDatabase: postgresDbConfig,
		S3:               s3Config,
	}
	return config
}
