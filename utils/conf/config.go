package conf

import (
	"time"

	"github.com/spf13/viper"
)

// Config menyimpan semua konfigurasi dari aplikasi
// nilai nya berasal dari viper yang membaca dari config file atau environtment variable
type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	MigrationURL        string        `mapstructure:"MIGRATION_URL"`
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
