package config

import (
	"reflect"
	"testing"
)

func Test_loadConfig(t *testing.T) {
	type args struct{}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "成功ケース: 環境変数が設定されていない場合、デフォルト値が設定される",
			args: args{},
			want: &Config{
				Env:                       "dev",
				Port:                      "8080",
				Database_url:              "",
				Sentry_DSN:                "",
				ProjectID:                 "",
				CLOUDFLARE_R2_ACCOUNT_ID:  "",
				CLOUDFLARE_R2_ACCESSKEY:   "",
				CLOUDFLARE_R2_SECRETKEY:   "",
				CLOUDFLARE_R2_BUCKET_NAME: "",
				COOKIE_DOMAIN:             "localhost",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
