package fileio

import (
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"
)

// Parse parses the file in the filesystem and returns the interface which
// can be casted into any type or struct as needed.
func (f *Provider) Parse(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		logrus.WithError(err).Error("Failed to read config file")
		return nil, err
	}

	var outputMap map[string]interface{}
	err = json.Unmarshal(data, &outputMap)
	if err != nil {
		logrus.WithError(err).Error("Failed to unmarshal config file")
		return nil, err
	}

	return outputMap, nil
}

// Save writes the interface to the filesystem.
func (f *Provider) Save(path string, data map[string]interface{}) error {
	jsonFile, err := json.Marshal(data)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal config")
		return err
	}

	err = os.WriteFile(path, jsonFile, os.ModePerm)
	if err != nil {
		logrus.WithError(err).Error("Failed to write config file")
		return err
	}

	return nil
}
