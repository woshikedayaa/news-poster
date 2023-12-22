package futil

import (
	"os"
	"path/filepath"
)

func GetRootDir() string {
	dir, _ := os.Getwd()
	return filepath.Dir(filepath.Dir(dir))
}
