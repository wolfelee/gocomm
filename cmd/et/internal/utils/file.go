package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// NL defines a new line
const (
	NL    = "\n"
	etDir = ".et"
)

// GetHome returns the path value of the et home where Join $HOME with .et
func GetHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, etDir), nil
}

// GetTemplateDir returns the category path value in etHome where could get it by GetEtHome
func GetTemplateDir(category string) (string, error) {
	etHome, err := GetHome()
	if err != nil {
		return "", err
	}

	return filepath.Join(etHome, category), nil
}

// FileExists returns true if the specified file is exists
func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

// LoadTemplate gets template content by the specified file
func LoadTemplate(category, file, builtin string) (string, error) {
	dir, err := GetTemplateDir(category)
	if err != nil {
		return "", err
	}

	file = filepath.Join(dir, file)
	if !FileExists(file) {
		return builtin, nil
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
