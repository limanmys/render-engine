package helpers

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// ReadConfiguration Retrieve data from liman
func ReadConfiguration() {
	customPath := os.Getenv("LIMAN_CONFIG")
	if customPath != "" {
		ConfigFilePath = customPath
	}

	viper.SetConfigFile(ConfigFilePath)
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Konfigürasyon dosyası okunamadı!\n%v\n", err.Error())
		os.Exit(1)
		return
	}

	// Set All Values to Variables
	AppKey = viper.GetString("APP_KEY")

	viper.SetDefault("DB_HOST", "127.0.0.1")
	DBHost = viper.GetString("DB_HOST")

	viper.SetDefault("DB_PORT", "5432")
	DBPort = viper.GetString("DB_PORT")

	viper.SetDefault("DB_DATABASE", "liman")
	DBName = viper.GetString("DB_DATABASE")

	viper.SetDefault("DB_USERNAME", "liman")
	DBUsername = viper.GetString("DB_USERNAME")
	DBPassword = viper.GetString("DB_PASSWORD")

	viper.SetDefault("LIMAN_RESTRICTED", false)
	RestrictedMode = viper.GetBool("LIMAN_RESTRICTED")

	viper.SetDefault("INTERNAL_ONLY", true)
	ListenInternally = viper.GetBool("INTERNAL_ONLY")

	viper.SetDefault("LOG_EXTENSION_PATH", "/liman/logs/extension.log")
	ExtensionLogsPath = viper.GetString("LOG_EXTENSION_PATH")

	viper.SetDefault("LOG_PATH", "/liman/logs/liman.log")
	LogsPath = viper.GetString("LOG_PATH")

	viper.SetDefault("SANDBOX_PATH", "/liman/sandbox/")
	SandboxPath = viper.GetString("SANDBOX_PATH")

	viper.SetDefault("KEYS_PATH", "/liman/keys/")
	KeysPath = viper.GetString("KEYS_PATH")

	viper.SetDefault("CERTS_PATH", "/liman/certs/")
	CertsPath = viper.GetString("CERTS_PATH")

	viper.SetDefault("EXTENSIONS_PATH", "/liman/extensions/")
	ExtensionsPath = viper.GetString("EXTENSIONS_PATH")

	CurrentIP = viper.GetString("CURRENT_IP")
}

func CheckRestrictedMode() bool {
	return RestrictedMode
}
