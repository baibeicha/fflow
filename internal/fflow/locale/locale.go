package locale

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

//go:embed messages/*.yaml
var messagesFS embed.FS

const appName = "fflow"

type Messages struct {
	CLI struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Long        string `yaml:"long"`
	} `yaml:"cli"`
	Commands map[string]struct {
		Short    string   `yaml:"short"`
		Long     string   `yaml:"long"`
		Examples []string `yaml:"examples"`
	} `yaml:"commands"`
	Flags    map[string]string `yaml:"flags"`
	Messages struct {
		Success  map[string]string `yaml:"success"`
		Errors   map[string]string `yaml:"errors"`
		Progress map[string]string `yaml:"progress"`
		Labels   map[string]string `yaml:"labels"`
		Results  map[string]string `yaml:"results"`
		Prompts  map[string]string `yaml:"prompts"`
	} `yaml:"messages"`
}

type LocaleConfig struct {
	Locale string `yaml:"locale"`
}

var (
	currentMessages *Messages
	currentLocale   string
	mu              sync.RWMutex
)

func init() {
	err := Init()
	if err != nil {
		slog.Error(err.Error())
	}
}

func Init() error {
	locale := detectLocale()
	return SetLocale(locale)
}

func detectLocale() string {
	configPath := getConfigPath()
	if configPath == "" {
		return "en"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return "en"
	}

	var cfg LocaleConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return "en"
	}

	if cfg.Locale == "" {
		return "en"
	}

	return cfg.Locale
}

func getConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(configDir, appName)
}

func getConfigPath() string {
	dir := getConfigDir()
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, "locale.yaml")
}

func SetLocale(locale string) error {
	locale = strings.ToLower(strings.TrimSpace(locale))
	if locale != "en" && locale != "ru" {
		return fmt.Errorf("invalid locale: %s", locale)
	}

	mu.Lock()
	defer mu.Unlock()

	msgs, err := loadMessages(locale)
	if err != nil {
		return err
	}

	currentMessages = msgs
	currentLocale = locale
	return nil
}

func loadMessages(locale string) (*Messages, error) {
	data, err := messagesFS.ReadFile(fmt.Sprintf("messages/%s.yaml", locale))
	if err != nil {
		return nil, fmt.Errorf("failed to load messages for locale %s: %w", locale, err)
	}

	var msgs Messages
	if err := yaml.Unmarshal(data, &msgs); err != nil {
		return nil, fmt.Errorf("failed to parse messages for locale %s: %w", locale, err)
	}

	return &msgs, nil
}

func GetLocale() string {
	mu.RLock()
	defer mu.RUnlock()
	return currentLocale
}

func GetMessages() *Messages {
	mu.RLock()
	defer mu.RUnlock()
	return currentMessages
}

func T(key string) string {
	mu.RLock()
	defer mu.RUnlock()

	if currentMessages == nil {
		return key
	}

	parts := strings.Split(key, ".")
	return getNestedValue(parts, currentMessages)
}

func Tf(key string, args ...interface{}) string {
	tmpl := T(key)
	if len(args) == 0 {
		return tmpl
	}
	return fmt.Sprintf(tmpl, args...)
}

func getNestedValue(parts []string, msgs *Messages) string {
	if len(parts) == 0 {
		return ""
	}

	switch parts[0] {
	case "cli":
		if len(parts) < 2 {
			return ""
		}
		switch parts[1] {
		case "name":
			return msgs.CLI.Name
		case "description":
			return msgs.CLI.Description
		case "long":
			return msgs.CLI.Long
		}
	case "commands":
		if len(parts) < 3 {
			return ""
		}
		cmd, ok := msgs.Commands[parts[1]]
		if !ok {
			return ""
		}
		switch parts[2] {
		case "short":
			return cmd.Short
		case "long":
			return cmd.Long
		}
	case "flags":
		if len(parts) < 2 {
			return ""
		}
		if val, ok := msgs.Flags[parts[1]]; ok {
			return val
		}
	case "messages":
		if len(parts) < 3 {
			return ""
		}
		switch parts[1] {
		case "success":
			if val, ok := msgs.Messages.Success[parts[2]]; ok {
				return val
			}
		case "errors":
			if val, ok := msgs.Messages.Errors[parts[2]]; ok {
				return val
			}
		case "progress":
			if val, ok := msgs.Messages.Progress[parts[2]]; ok {
				return val
			}
		case "labels":
			if val, ok := msgs.Messages.Labels[parts[2]]; ok {
				return val
			}
		case "results":
			if val, ok := msgs.Messages.Results[parts[2]]; ok {
				return val
			}
		case "prompts":
			if val, ok := msgs.Messages.Prompts[parts[2]]; ok {
				return val
			}
		}
	}

	return parts[len(parts)-1]
}

func SaveLocale(locale string) error {
	locale = strings.ToLower(strings.TrimSpace(locale))
	if locale != "en" && locale != "ru" {
		return fmt.Errorf("invalid locale: %s", locale)
	}

	configDir := getConfigDir()
	if configDir == "" {
		return fmt.Errorf("cannot determine config directory")
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := getConfigPath()
	cfg := LocaleConfig{Locale: locale}
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
