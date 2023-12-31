package conf

import (
	"time"

	"github.com/spf13/viper"
)

// Config menyimpan semua konfigurasi dari aplikasi
// nilai nya berasal dari viper yang membaca dari config file atau environtment variable
type Config struct {
	DBSource             string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	Environment          string        `mapstructure:"ENVIRONMENT"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	EmailSenderName      string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress   string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword  string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

// LoadConfig membaca configuration dari file jika ada atau environtment variables jika disediakan
func LoadConfig(path string) (config Config, err error) {
	// SETUP VIPER UNTUK MEMBACA DARI FILE
	viper.AddConfigPath(path)  // memberitahu viper lokasi configuration file nya
	viper.SetConfigName("app") // memberitahu viper untuk mencari file dengan nama yang  ada di argument "app" dari app.env
	viper.SetConfigType("env") // memberitahu tipe filenya. bisa JSON, XML dll

	// SETUP VIPER UNTUK MEMBACA DARI ENVIRONTMENT. JIKA ADA MAKA OVERRIDE SETUP DIATAS
	viper.AutomaticEnv() // otomatis mengoverride nilai dari file configuration jika ada configuration dari environment

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config) // unmarshal the values into the target config object.
	return
}
