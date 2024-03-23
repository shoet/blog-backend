package config

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	Env                         string `env:"BLOG_ENV,required"`
	AppPort                     int64  `env:"BLOG_APP_PORT,required"`
	LogLevel                    string `env:"BLOG_LOG_LEVEL" envDefault:"info"`
	DBHost                      string `env:"BLOG_DB_HOST,required"`
	DBPort                      int64  `env:"BLOG_DB_PORT,required"`
	DBUser                      string `env:"BLOG_DB_USER,required"`
	DBPass                      string `env:"BLOG_DB_PASS,required"`
	DBName                      string `env:"BLOG_DB_NAME,required"`
	DBTlsEnabled                bool   `env:"BLOG_DB_TLS_ENABLED" envDefault:"false"`
	DBSSLMode                   string `env:"BLOG_DB_SSL_MODE" envDefault:"disable"`
	KVSHost                     string `env:"BLOG_KVS_HOST,required"`
	KVSPort                     int64  `env:"BLOG_KVS_PORT,required"`
	KVSUser                     string `env:"BLOG_KVS_USER,required"`
	KVSPass                     string `env:"BLOG_KVS_PASS,required"`
	KVSTlsEnabled               bool   `env:"BLOG_KVS_TLS_ENABLED" envDefault:"false"`
	AWSS3Region                 string `env:"AWS_DEFAULT_REGION"`
	AWSS3Bucket                 string `env:"BLOG_AWS_S3_BUCKET,required"`
	AWSS3ThumbnailDirectory     string `env:"BLOG_AWS_S3_THUMBNAIL_DIRECTORY,required"`
	AWSSS3ContentImageDirectory string `env:"BLOG_AWS_S3_CONTENT_IMAGE_DIRECTORY,required"`
	AWSS3PresignPutExpiresSec   int64  `env:"BLOG_AWS_S3_PRESIGN_PUT_EXPIRES_SEC" envDefault:"300"`
	AdminName                   string `env:"ADMIN_NAME,required"`
	AdminEmail                  string `env:"ADMIN_EMAIL,required"`
	AdminPassword               string `env:"ADMIN_PASSWORD,required"`
	JWTSecret                   string `env:"JWT_SECRET,required"`
	JWTExpiresInSec             int    `env:"JWT_EXPIRES_IN_SEC" envDefault:"86400"`
	CORSWhiteList               string `env:"CORS_WHITE_LIST"`
	SiteDomain                  string `env:"SITE_DOMAIN"`
	CdnDomain                   string `env:"CDN_DOMAIN"`
	GitHubPersonalAccessToken   string `env:"GITHUB_PERSONAL_ACCESS_TOKEN"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
