package config

import (
	"context"
	"os"
	"reflect"
)

type Env string

const CtxEnvKey Env = "env"

type Config struct {
	Env                       string `env:"ENV" envDefault:"dev"`
	Port                      string `env:"PORT" envDefault:"8080"`
	Database_url              string `env:"DATABASE_URL" envDefult:""`
	Sentry_DSN                string `env:"SENTRY_DSN" envDefult:""`
	ProjectID                 string `env:"PROJECT_ID" envDefult:""`
	CLOUDFLARE_R2_ACCOUNT_ID  string `env:"CLOUDFLARE_R2_ACCOUNT_ID" envDefult:""`
	CLOUDFLARE_R2_ACCESSKEY   string `env:"CLOUDFLARE_R2_ACCESSKEY" envDefult:""`
	CLOUDFLARE_R2_SECRETKEY   string `env:"CLOUDFLARE_R2_SECRETKEY" envDefult:""`
	CLOUDFLARE_R2_BUCKET_NAME string `env:"CLOUDFLARE_R2_BUCKET_NAME" envDefult:""`
	COOKIE_DOMAIN             string `env:"COOKIE_DOMAIN" envDefault:"localhost"`
}

func New(ctx context.Context) (context.Context, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, CtxEnvKey, cfg), nil
}

func loadConfig() (*Config, error) {
	cfg := &Config{}
	value := reflect.Indirect(reflect.ValueOf(cfg))

	for _, v := range reflect.VisibleFields(value.Type()) {
		env := os.Getenv(v.Tag.Get("env"))
		if env != "" {
			value.FieldByName(v.Name).SetString(env)
			continue
		}
		defaultValue := v.Tag.Get("envDefault")
		value.FieldByName(v.Name).SetString(defaultValue)
	}

	return cfg, nil
}
