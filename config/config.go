/*
Package config handles the loading and validation of the database checker configuration.
It supports YAML-based configuration for multiple database instances and their connection details.
*/
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DatabaseConfig defines the connection and health check parameters for a single database.
type DatabaseConfig struct {
	Type           string `yaml:"type"`
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	Name           string `yaml:"name"`
	HealthQuery    string `yaml:"health_query"`
	TLSMode        string `yaml:"tls_mode,omitempty"`
	WalletPath     string `yaml:"wallet_path,omitempty"`
	RootCertPath   string `yaml:"root_cert_path,omitempty"`
	ClientCertPath string `yaml:"client_cert_path,omitempty"`
	ClientKeyPath  string `yaml:"client_key_path,omitempty"`
}

// Config holds the collection of database configurations indexed by a unique identifier.
type Config struct {
	Databases map[string]DatabaseConfig `yaml:"databases"`
}

var supportedDBTypes = map[string]struct{}{"mysql": {}, "postgres": {}, "oracle": {}, "sqlserver": {}, "sqlite": {}, "mongodb": {}}
var supportedTLSModes = map[string]struct{}{"disable": {}, "require": {}, "verify-ca": {}, "verify-full": {}, "": {}}

// Validate checks the configuration for any unsupported or invalid values.
func (c *Config) Validate() error {
	for id, dbConfig := range c.Databases {
		if _, ok := supportedDBTypes[dbConfig.Type]; !ok {
			return fmt.Errorf("database %q has unsupported type: %s", id, dbConfig.Type)
		}
		if _, ok := supportedTLSModes[dbConfig.TLSMode]; !ok {
			return fmt.Errorf("database %q has unsupported tls_mode: %s", id, dbConfig.TLSMode)
		}
	}
	return nil
}

// LoadConfig reads a YAML configuration file from disk.
// It uses os.OpenRoot (Go 1.24+) to safely access the file and prevent directory traversal.
func LoadConfig(configFile string) (*Config, error) {
	root, err := os.OpenRoot(filepath.Dir(configFile))
	if err != nil {
		return nil, fmt.Errorf("failed to open config directory: %w", err)
	}
	defer func() {
		_ = root.Close()
	}()

	data, err := root.ReadFile(filepath.Base(configFile))
	if err != nil {
		return nil, err
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err = config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &config, nil
}
