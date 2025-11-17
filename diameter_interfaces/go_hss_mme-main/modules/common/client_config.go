package common

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// MMEConfig represents the configuration of a single MME client
type MMEConfig struct {
	IDMmeIdentity  int    `yaml:"idmmeidentity"`
	MMEHost        string `yaml:"mmehost"`
	MMERealm       string `yaml:"mmerealm"`
	UEReachability int    `yaml:"UE_reachability"`
	IP             string `yaml:"IP"`
}

// Wrapper struct to handle the 'clients' key in the YAML file
type ClientsConfig struct {
	Clients []MMEConfig `yaml:"clients"`
}

// Global variable to hold the loaded MME configuration (multiple clients)
var mmeInstances []MMEConfig

// LoadMMEConfig loads the MME configuration from a YAML file
func LoadMMEConfig(configPath string) error {
	// Check if the file exists
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("MME config file does not exist: %s", configPath)
	}

	// Read the file
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading MME config file: %v", err)
	}

	// Parse the YAML content into the wrapper struct
	var config ClientsConfig
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return fmt.Errorf("error parsing MME config file: %v", err)
	}

	// Validate that there are clients in the configuration
	if len(config.Clients) == 0 {
		return errors.New("MME config is missing client entries")
	}

	// Set the global instance to the parsed clients list
	mmeInstances = config.Clients
	return nil
}

// GetMMEConfig returns the loaded MME configuration for all clients
func GetMMEConfig() ([]MMEConfig, error) {
	if len(mmeInstances) == 0 {
		return nil, errors.New("MME config is not loaded. Please call LoadMMEConfig first")
	}
	return mmeInstances, nil
}
