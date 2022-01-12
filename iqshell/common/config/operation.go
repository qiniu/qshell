package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func NewConfigWithPath(path string) (config *Config, err error) {
	file, err := os.Open(path)
	if err != nil || file == nil {
		return
	}

	defer func(file *os.File) {
		if e := file.Close(); e != nil && err != nil {
			err = e
		}
	}(file)

	configData, err := ioutil.ReadAll(file)
	if err != nil || configData == nil {
		return
	}

	err = json.Unmarshal(configData, &config)
	return
}

func (c *Config) UpdateToLocal(path string) (err error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0666)
	if err != nil || file == nil {
		return
	}

	defer func(file *os.File) {
		if e := file.Close(); e != nil && err != nil {
			err = e
		}
	}(file)

	configData, err := json.MarshalIndent(c, "", "\t")
	if err != nil || configData == nil || len(configData) == 0 {
		return
	}

	fmt.Println("configData:" + string(configData))
	writer := bufio.NewWriter(file)

	_, err = writer.Write(configData)
	if err != nil {
		return
	}

	err = writer.Flush()

	return
}