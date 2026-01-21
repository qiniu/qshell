package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

func NewConfigWithPath(path string) (*Config, *data.CodeError) {
	file, err := os.Open(path)
	if err != nil || file == nil {
		return nil, data.NewEmptyError().AppendError(err)
	}

	defer func(file *os.File) {
		file.Close()
	}(file)

	configData, err := ioutil.ReadAll(file)
	if err != nil || configData == nil {
		return nil, data.NewEmptyError().AppendError(err)
	}

	config := &Config{}
	if e := json.Unmarshal(configData, &config); e != nil {
		return nil, data.NewEmptyError().AppendError(err)
	}
	return config, nil
}

func (c *Config) UpdateToLocal(path string) *data.CodeError {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0o666)
	if err != nil || file == nil {
		return data.NewEmptyError().AppendError(err)
	}

	defer func(file *os.File) {
		file.Close()
	}(file)

	configData, err := json.MarshalIndent(c, "", "\t")
	if err != nil || configData == nil || len(configData) == 0 {
		return data.ConvertError(err)
	}

	fmt.Println("configData:" + string(configData))
	writer := bufio.NewWriter(file)

	_, err = writer.Write(configData)
	if err != nil {
		return data.NewEmptyError().AppendError(err)
	}

	err = writer.Flush()

	return data.ConvertError(err)
}
