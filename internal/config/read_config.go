package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	DB_URL            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() (Config, error) {
	// Getting the correct filepath for gatorconfig.json
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("ERROR getting config filepath: %v", err)
	}

	// Checking if gatorconfig.json exists and creating if it doesn't
	_, e := os.Stat(configFilePath)
	if e != nil {
		if errors.Is(e, os.ErrNotExist) {
			err = write(Config{DB_URL: "postgres://example"})
			if err != nil {
				return Config{}, err
			}
		} else {
			// Panic if another error occurs outside of ErrNotExist
			return Config{}, fmt.Errorf("ERROR reading config file: %v", e)
		}
	}

	// Decoding config file to struct
	body, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("ERROR opening config file: %v", err)
	}
	defer body.Close()

	var cfg Config
	decoder := json.NewDecoder(body)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("ERROR decoding config file to struct: %v", err)
	}

	return cfg, nil
}

func (cfg Config) SetUser(uname string) error {
	cfg.Current_user_name = uname
	err := write(cfg)
	if err != nil {
		return err
	}
	return nil
}

// Helper functions
const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/" + configFileName, nil
}

func write(cfg Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("ERROR getting config filepath: %v", err)
	}

	file, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("ERROR openning config file: %v", err)
	}
	defer file.Close()

	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("ERROR marshalling config struct data: %v", err)
	}
	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("ERROR writing config struct to file: %v", err)
	}
	err = file.Sync()
	if err != nil {
		return fmt.Errorf("ERROR saving over config file: %v", err)
	}
	return nil
}
