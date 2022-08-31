package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var ErrConfigDoesNotExist = errors.New("Config file does not exist")

type (
	RdmConfig struct {
		Commands map[string]*UserCommand `json:"commands"`
	}

	UserCommand struct {
		ExecutablePath string `json:"executablePath"`
	}
)

func New() *RdmConfig {
	return &RdmConfig{Commands: map[string]*UserCommand{}}
}

func (r *RdmConfig) Load(configPath string) error {
	_, err := os.Stat(configPath)

	if err != nil {
		return ErrConfigDoesNotExist
	}

	contents, err := ioutil.ReadFile(configPath)

	if err != nil {
		return fmt.Errorf("Could not load config: %w", err)
	}

	err = json.Unmarshal(contents, &r)
	if err != nil {
		return fmt.Errorf("Could not load config: %w", err)
	}

	relativeRoot := filepath.Dir(configPath)

	for _, command := range r.Commands {
		if !filepath.IsAbs(command.ExecutablePath) {
			path := filepath.Join(relativeRoot, command.ExecutablePath)
			command.ExecutablePath = path
		}
	}

	return nil
}
