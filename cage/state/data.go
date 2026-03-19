package state

import (
	"cage/cage/config"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

var (
	DataDirPath = filepath.Join(xdg.DataHome, config.AppName)
	DataDir     *os.Root
)

type CorruptDataDirErr struct {
	Path string
}

func (c CorruptDataDirErr) Error() string {
	return fmt.Sprintf("data directory is corrupt: %v", c.Path)
}

func IsInitialized() bool {
	dirPath := DataDirPath
	_, exists, err := LoadDataDir(dirPath)
	if err != nil {
		return false
	}
	if !exists {
		return false
	}
	return true
}

func Init() error {
	dirPath := DataDirPath
	dir, exists, err := LoadDataDir(dirPath)
	if err != nil {
		return err
	}
	if !exists {
		err = os.MkdirAll(dirPath, 0760)
		if err != nil {
			return err
		}
		dir, err = os.OpenRoot(dirPath)
		if err != nil {
			return err
		}
	}

	DataDir = dir

	return nil
}

func LoadDataDir(dirPath string) (dir *os.Root, exists bool, err error) {
	fi, err := os.Stat(dirPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	if !fi.IsDir() {
		return nil, false, CorruptDataDirErr{Path: fi.Name()}
	}

	dir, err = os.OpenRoot(dirPath)
	if err != nil {
		return nil, false, err
	}

	err = InitNixStore()
	if err != nil {
		return nil, false, err
	}

	return dir, true, nil
}

func init() {
	dir, exists, err := LoadDataDir(DataDirPath)
	if err != nil {
		return
	}
	if !exists {
		return
	}
	DataDir = dir
}
