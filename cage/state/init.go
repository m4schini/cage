package state

import (
	"cage/cage/config"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var CageAlreadyExistsErr = fmt.Errorf("cage already exists")

func NewCage() (cfg CageConfig, err error) {
	fi, err := os.Stat("..")
	if err != nil {
		return cfg, err
	}

	cfg.Name = fi.Name()
	cfg.Env = append(cfg.Env, EnvVar{
		Key:   "CAGE",
		Value: "",
	})
	return cfg, nil
}

func CreateCageDir(cfg CageConfig) error {
	cagePath := "./" + config.CageDirName()
	_, err := os.Stat(cagePath)
	if !errors.Is(err, os.ErrNotExist) {
		return CageAlreadyExistsErr
	}

	err = os.MkdirAll(cagePath, 0750)
	if err != nil {
		return err
	}
	configFilePath := filepath.Join(".", config.CageDirName(), "config.yaml")

	f, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = yaml.NewEncoder(f).Encode(cfg)
	return err
}

func CageExists() bool {
	cagePath := "./" + config.CageDirName()
	_, err := os.Stat(cagePath)
	return !errors.Is(err, os.ErrNotExist)
}
