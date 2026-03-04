package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	AppID     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
	Region    string `mapstructure:"region"`
	Defaults  struct {
		Timezone        string `mapstructure:"timezone"`
		ReminderMinutes int    `mapstructure:"reminder_minutes"`
	} `mapstructure:"defaults"`
	OAuth struct {
		RedirectPort int `mapstructure:"redirect_port"`
	} `mapstructure:"oauth"`
	CustomEmojis map[string]string `mapstructure:"custom_emojis"`
}

var (
	cfg     *Config
	cfgDir  string
	rootDir string
)

// GetConfigDir returns the .lark directory path
func GetConfigDir() string {
	return cfgDir
}

// GetRootDir returns the project root directory
func GetRootDir() string {
	return rootDir
}

// expandPath expands a leading ~ to the user's home directory
func expandPath(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, path[1:]), nil
}

// Init initializes the configuration
func Init() error {
	// Config directory can be set via LARK_CONFIG_DIR or legacy LARK_CAL_CONFIG_DIR
	cfgDir = os.Getenv("LARK_CONFIG_DIR")
	if cfgDir == "" {
		cfgDir = os.Getenv("LARK_CAL_CONFIG_DIR") // Legacy fallback
	}
	if cfgDir == "" {
		return fmt.Errorf("LARK_CONFIG_DIR environment variable is not set")
	}

	var err error
	cfgDir, err = expandPath(cfgDir)
	if err != nil {
		return err
	}

	rootDir = filepath.Dir(cfgDir)

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(cfgDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(cfgDir)

	// Set defaults
	viper.SetDefault("defaults.timezone", "Asia/Singapore")
	viper.SetDefault("defaults.reminder_minutes", 15)
	viper.SetDefault("oauth.redirect_port", 9999)
	viper.SetDefault("region", "lark")

	// Environment variable bindings
	viper.SetEnvPrefix("LARK")
	viper.BindEnv("app_id", "LARK_APP_ID")
	viper.BindEnv("app_secret", "LARK_APP_SECRET")
	viper.BindEnv("region", "LARK_REGION")

	// Read config file (if exists)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config: %w", err)
		}
		// Config file not found is OK, we'll use defaults and env vars
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// Get returns the current configuration
func Get() *Config {
	if cfg == nil {
		cfg = &Config{}
	}
	return cfg
}

// GetAppID returns the app ID from config or environment
func GetAppID() string {
	return viper.GetString("app_id")
}

// GetAppSecret returns the app secret from environment
func GetAppSecret() string {
	return viper.GetString("app_secret")
}

// GetTimezone returns the default timezone
func GetTimezone() string {
	return viper.GetString("defaults.timezone")
}

// GetRedirectPort returns the OAuth redirect port
func GetRedirectPort() int {
	return viper.GetInt("oauth.redirect_port")
}

// TokensFilePath returns the path to the tokens file
func TokensFilePath() string {
	return filepath.Join(cfgDir, "tokens.json")
}

// TenantTokensFilePath returns the path to the tenant tokens file
func TenantTokensFilePath() string {
	return filepath.Join(cfgDir, "tenant_tokens.json")
}

// GetCustomEmojis returns the custom emoji mappings
func GetCustomEmojis() map[string]string {
	return viper.GetStringMapString("custom_emojis")
}

// GetRegion returns the configured region ("lark" or "feishu")
func GetRegion() string {
	return viper.GetString("region")
}

// GetBaseURL returns the API base URL for the configured region
func GetBaseURL() string {
	if GetRegion() == "feishu" {
		return "https://open.feishu.cn/open-apis"
	}
	return "https://open.larksuite.com/open-apis"
}
