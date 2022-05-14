package fileio

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// CheckFile checks if a file exists.
func (f *Provider) CheckFile(path string) bool {
	_, err := os.ReadFile(path)
	return err == nil
}

// GenerateFile generates a file.
func (f *Provider) GenerateFile(path string) error {

	// get the path until the last backslash
	dir := path[:len(path)-len(filepath.Base(path))]
	// check if the directory exists
	if !f.CheckFolder(dir) {
		// create the directory
		err := f.GenerateFolder(dir)
		if err != nil {
			logrus.WithError(err).Error("Failed to generate folder")
			return err
		}
	}
	_, err := os.Create(path)
	if err != nil {
		logrus.WithError(err).Error("Failed to create file")
		return err
	}
	return nil
}

// DeleteFile deletes a file.
func (f *Provider) DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete file")
		return err
	}
	return nil
}

// CheckFolder checks if a folder exists.
func (f *Provider) CheckFolder(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GenerateFolder generates a folder.
func (f *Provider) GenerateFolder(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate folder")
		return err
	}
	return nil
}

// DeleteFolder deletes a folder.
func (f *Provider) DeleteFolder(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete folder")
		return err
	}
	return nil
}
