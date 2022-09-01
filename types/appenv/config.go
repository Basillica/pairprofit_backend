package appenv

type AppConfig struct {
	AWS_REGION            string `mapstructure:"AWS_REGION"`
	COOKIE_DOMAIN         string `mapstructure:"COOKIE_DOMAIN"`
	APP_ENVIRONMENT       string `mapstructure:"APP_ENVIRONMENT"`
	COOKIE_SECURE_ENABLE  bool   `mapstructure:"COOKIE_SECURE_ENABLE"`
	DECODING_SECRET       string `mapstructure:"DECODING_SECRET"`
	ALLOWED_ORIGIN        string `mapstructure:"ALLOWED_ORIGIN"`
	COOKIE_HTTPONLY       bool   `mapstructure:"COOKIE_HTTPONLY"`
	HUBSPOT_ENABLE        bool   `mapstructure:"HUBSPOT_ENABLE"`
	CRYPTO_SECRET         string `mapstructure:"CRYPTO_SECRET"`
	COGNITO_POOL_ID       string `mapstructure:"COGNITO_POOL_ID"`
	COGNITO_CLIENT_ID     string `mapstructure:"COGNITO_CLIENT_ID"`
	COGNITO_CLIENT_SECRET string `mapstructure:"COGNITO_CLIENT_SECRET"`
	SIB_API_KEY           string `mapstructure:"SIB_API_KEY"`
}
