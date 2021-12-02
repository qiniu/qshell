package utils

import "github.com/mitchellh/go-homedir"

func GetHomePath() (string, error) {
	return homedir.Dir()
}

func CreateDirIfNotExist(path string) error {
	return nil
}