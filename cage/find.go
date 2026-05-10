package cage

import (
	"cage/cage/state"
	"errors"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

var ErrNotFound = errors.New(".cage directory not found")

func FindCageDir(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		candidate := filepath.Join(dir, ".cage")
		fi, err := os.Stat(candidate)
		if err == nil && fi.IsDir() {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ErrNotFound
		}
		dir = parent
	}
}

func LoadConfig() (cfg state.CageConfig, lastModified time.Time, err error) {
	cageDir, err := FindCageDir("./")
	if err != nil {
		return state.CageConfig{}, time.Now(), err
	}

	path := filepath.Join(cageDir, "config.yaml")

	fi, err := os.Stat(path)
	if err != nil {
		return state.CageConfig{}, time.Now(), err
	}

	f, err := os.Open(path)
	if err != nil {
		return state.CageConfig{}, time.Now(), err
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&cfg)
	return cfg, fi.ModTime(), err
}
