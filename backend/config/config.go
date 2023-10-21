package config

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	Env                       string `env:"BLOG_ENV"`
	AppPort                   int64  `env:"BLOG_APP_PORT"`
	DBHost                    string `env:"BLOG_DB_HOST"`
	DBPort                    int64  `env:"BLOG_DB_PORT"`
	DBUser                    string `env:"BLOG_DB_USER"`
	DBPass                    string `env:"BLOG_DB_PASS"`
	DBName                    string `env:"BLOG_DB_NAME"`
	KVSHost                   string `env:"BLOG_KVS_HOST"`
	KVSPort                   int64  `env:"BLOG_KVS_PORT"`
	KVSUser                   string `env:"BLOG_KVS_USER"`
	KVSPass                   string `env:"BLOG_KVS_PASS"`
	AWSS3Region               string `env:"AWS_DEFAULT_REGION"`
	AWSS3Bucket               string `env:"BLOG_AWS_S3_BUCKET"`
	AWSS3ThumbnailDirectory   string `env:"BLOG_AWS_S3_THUMBNAIL_DIRECTORY"`
	AWSS3PresignPutExpiresSec int64  `env:"BLOG_AWS_S3_PRESIGN_PUT_EXPIRES_SEC" envDefault:"300"`
	AdminEmail                string `env:"ADMIN_EMAIL"`
	AdminPassword             string `env:"ADMIN_PASSWORD"`
	JWTSecret                 string `env:"JWT_SECRET"`
	JWTExpiresInSec           int    `env:"JWT_EXPIRES_IN_SEC" envDefault:"86400"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
