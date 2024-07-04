package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the layout of the configuration file.
type Config struct {
	Title       string             `yaml:"site_title"`
	URL         string             `yaml:"site_url"`
	Description string             `yaml:"description"`
	Theme       string             `yaml:"theme"`
	Syntax      SyntaxHighlighting `yaml:"syntax_highlighting"`
	Nav         []NavItem          `yaml:"nav"`
	Social      []NavItem          `yaml:"social"`
	Port        int                `yaml:"port"`
}

// SyntaxHighlighting contains settings for syntax highlighting themes.
type SyntaxHighlighting struct {
	DarkMode  ThemeConfig `yaml:"dark_mode"`
	LightMode ThemeConfig `yaml:"light_mode"`
}

// ThemeConfig holds the theme settings.
type ThemeConfig struct {
	Theme string `yaml:"theme"`
}

// NavItem represents a navigation item.
type NavItem struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// LoadConfig attempts to load a configuration from a "config.yml" file
// in the current directory. If the file does not exist, it falls back
// to an environment variable to retrieve the config path.
func LoadConfig() (*Config, error) {
	path := "config.yml"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		path = os.Getenv("SITE_CONFIG")
		if path == "" {
			return nil, fmt.Errorf("config file not found and SITE_CONFIG environment variable is not set")
		}
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	cfg.setDefaults()

	// Validate the loaded config
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid cfg: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets required default values for the configuration
func (c *Config) setDefaults() {
	if c.Title == "" {
		c.Title = "some title"
	}
	if c.URL == "" {
		c.URL = "http://localhost"
	}
	if c.Port == 0 {
		c.Port = 8080
	}
	if c.Theme == "" {
		c.Theme = "default"
	}
	if c.Syntax.DarkMode.Theme == "" {
		c.Syntax.DarkMode.Theme = "monokai"
	}
	if c.Syntax.LightMode.Theme == "" {
		c.Syntax.LightMode.Theme = "github"
	}
}

func (c *Config) validate() error {
	var errors []string

	if c.URL == "" {
		errors = append(errors, "site URL is required")
	}
	if c.Port <= 0 {
		errors = append(errors, "port must be a positive integer")
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil
}
