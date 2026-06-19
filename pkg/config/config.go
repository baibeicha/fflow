package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
	filename string
	Log      LogConfig `mapstructure:"log"`
}

type LogConfig struct {
	Type  string `mapstructure:"type"`
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

func init() {
	_ = godotenv.Load()
}

func MustLoadEmpty(defaults ...string) *Config {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	cfg.Viper = v
	for i := 0; i < len(defaults); i += 2 {
		if !v.IsSet(defaults[i]) {
			v.SetDefault(defaults[i], defaults[i+1])
		}
	}

	return &cfg
}

func MustLoad(filename string, defaults ...string) *Config {
	if filename == "" {
		filename = "config"
	}

	v := viper.New()
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("./configs")

	for i := 0; i < len(defaults); i++ {
		if defaults[i] == "cfg-path" && i+1 < len(defaults) {
			v.AddConfigPath(defaults[i+1])
			defaults = slices.Delete(defaults, i, i+2)
			i--
		}
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error finding config file: %v", err)
	}

	configFile := v.ConfigFileUsed()

	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("error reading config file %s: %v", configFile, err)
	}

	expandedConfig := os.ExpandEnv(string(configBytes))

	configExt := strings.TrimPrefix(filepath.Ext(configFile), ".")
	v.SetConfigType(configExt)

	if err := v.ReadConfig(strings.NewReader(expandedConfig)); err != nil {
		log.Fatalf("error parsing expanded config: %v", err)
	}

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("error parsing config into struct: %v", err)
	}

	for i := 0; i < len(defaults); i += 2 {
		if !v.IsSet(defaults[i]) && i+1 < len(defaults) {
			v.SetDefault(defaults[i], defaults[i+1])
		}
	}

	cfg.Viper = v
	cfg.filename = filename

	return &cfg
}

func GetTimeUnit(unit string) time.Duration {
	switch unit {
	case "s":
		return time.Second
	case "m":
		return time.Minute
	case "h":
		return time.Hour
	case "d":
		return time.Hour * 24
	default:
		return time.Minute
	}
}

func MustLoadUserConfig(appName string, defaults ...string) *Config {
	v := viper.New()

	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}

	appDir := filepath.Join(configDir, appName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		log.Printf("warning: cannot create config directory %s: %v", appDir, err)
	}

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(appDir)
	v.AddConfigPath(".")

	v.SetEnvPrefix(strings.ToUpper(appName))
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	for i := 0; i < len(defaults); i += 2 {
		if i+1 < len(defaults) {
			v.SetDefault(defaults[i], defaults[i+1])
		}
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok || os.IsNotExist(err) {
			configPath := filepath.Join(appDir, "config.yaml")
			if writeErr := writeDefaultConfig(configPath, defaults); writeErr != nil {
				log.Printf("warning: cannot write default config: %v", writeErr)
			}
			_ = v.ReadInConfig()
		} else {
			log.Printf("warning: error reading user config: %v", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Printf("warning: error parsing user config: %v", err)
	}

	cfg.Viper = v
	cfg.filename = filepath.Join(appDir, "config")

	return &cfg
}

func writeDefaultConfig(path string, defaults []string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString("# fflow user configuration\n")
	sb.WriteString("# This file is auto-generated. You can modify it manually.\n\n")

	for i := 0; i < len(defaults); i += 2 {
		if i+1 < len(defaults) {
			key := defaults[i]
			value := defaults[i+1]
			sb.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		}
	}

	return os.WriteFile(path, []byte(sb.String()), 0644)
}
