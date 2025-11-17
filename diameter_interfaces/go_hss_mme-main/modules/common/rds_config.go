package common

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// RDSConfig holds the database settings loaded from rds.yaml
type RDSConfig struct {
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Pool struct {
		MaxOpenConns    int `yaml:"maxOpenConns"`
		MaxIdleConns    int `yaml:"maxIdleConns"`
		ConnMaxLifetime int `yaml:"connMaxLifetime"` // In seconds
	} `yaml:"pool"`
}

// Global variable to hold the loaded RDS configuration
var RDS RDSConfig

// LoadRDSConfig loads the RDS configuration from a YAML file
func LoadRDSConfig(configPath string) error {
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading RDS config file: %v", err)
	}

	if err := yaml.Unmarshal(configFile, &RDS); err != nil {
		return fmt.Errorf("error parsing RDS config file: %v", err)
	}

	return nil
}
