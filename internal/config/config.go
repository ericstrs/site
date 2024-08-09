package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var ErrConfigNotFound = errors.New("config file not found and SITE_CONFIG environment variable is not set")

// Config represents the layout of the configuration file.
type Config struct {
	Title       string             `yaml:"title"`
	URL         string             `yaml:url`
	Host        string             `yaml:"host"`
	Port        int                `yaml:"port"`
	Description string             `yaml:"description"`
	Theme       string             `yaml:"theme"`
	Syntax      SyntaxHighlighting `yaml:"syntax_highlighting"`
	Nav         []NavItem          `yaml:"nav"`
	Social      []NavItem          `yaml:"social"`
	DocsPath    string             `yaml:"docs_path"`
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

// LoadConfig loads or initializes the config file and ensures
// the "docs" directory exists. It returns a pointer to the Config
// struct and any error encountered during the process.
func LoadConfig() (*Config, error) {
	cfg, err := loadFile()
	if err != nil {
		if errors.Is(err, ErrConfigNotFound) {
			path := "config.yml" // Default to "config.yml" in the current directory
			if err := createConfigFile(path); err != nil {
				return nil, fmt.Errorf("failed to create config file: %w", err)
			}
			cfg, err = loadFile()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if cfg.DocsPath == "" {
		docsPath, err := findDocsDir()
		if err != nil {
			docsPath = "docs" // Default to "docs" in current directory
			if err := os.MkdirAll(docsPath, 0755); err != nil {
				return nil, fmt.Errorf("failed to resolve docs path and failed to create docs directory: %w", err)
			}

			readmePath := filepath.Join(docsPath, "README.md")
			if _, err := os.Create(readmePath); err != nil {
				return nil, fmt.Errorf("failed to create README.md: %w", err)
			}
		}

		cfg.DocsPath = docsPath
	}

	return cfg, nil
}

// loadFile attempts to read and unmarshal the configuration file.
func loadFile() (*Config, error) {
	path, err := findConfigDir()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return &cfg, nil
}

// findConfigDir first attempts to locate the configuration file in the
// current directory and then at SITE_CONFIG.
func findConfigDir() (string, error) {
	path := "config.yml"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		path = os.Getenv("SITE_CONFIG")
		if path == "" {
			return "", ErrConfigNotFound
		}
	}
	return path, nil
}

// createConfigFile creates a default configuration file.
func createConfigFile(path string) error {
	defaultConfig := &Config{
		Title:       "some title",
		URL:         "http://localhost:8080",
		Host:        "localhost",
		Port:        8080,
		Description: "",
		Theme:       "default",
		Syntax: SyntaxHighlighting{
			DarkMode:  ThemeConfig{Theme: "monokai"},
			LightMode: ThemeConfig{Theme: "github"},
		},
		Nav:    []NavItem{},
		Social: []NavItem{},
	}

	content, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %w", err)
	}

	if err := ioutil.WriteFile(path, content, 0644); err != nil {
		return fmt.Errorf("failed to write default config file: %w", err)
	}

	return nil
}

// findDocsDir attempts to locate the "docs" directory in the current directory.
// If the "docs" directory does not exist, it falls back on DOCS_DIR to
// retrieve the path to the docs directory.
func findDocsDir() (string, error) {
	path := "docs"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		path = os.Getenv("DOCS_DIR")
		if path == "" {
			return "", fmt.Errorf("docs directory not found and DOCS_DIR environment variable is not set")
		}
	}
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("failed to locate docs directory: %w", err)
	}
	return path, nil
}
