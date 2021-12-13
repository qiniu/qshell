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

func (c *Config) Merge(from *Config) {
	if from == nil {
		return
	}

	if !c.HasCredentials() {
		c.Credentials = from.Credentials
	}

	if len(c.UseHttps) == 0 {
		c.UseHttps = from.UseHttps
	}

	c.Hosts.merge(&from.Hosts)
	c.Up.merge(&from.Up)
	c.Download.merge(&from.Download)
}

func (c *Hosts) merge(from *Hosts) {
	if from == nil {
		return
	}

	if len(c.UC) == 0 {
		c.UC = from.UC
	}

	if len(c.Api) == 0 {
		c.Api = from.Api
	}

	if len(c.Rsf) == 0 {
		c.Rsf = from.Rsf
	}

	if len(c.Rs) == 0 {
		c.Rs = from.Rs
	}

	if len(c.Io) == 0 {
		c.Io = from.Io
	}

	if len(c.Up) == 0 {
		c.Up = from.Up
	}
}

func (up *Up) merge(from *Up) {
	if from == nil {
		return
	}

	if up.PutThreshold == 0 {
		up.PutThreshold = from.PutThreshold
	}

	if up.ChunkSize == 0 {
		up.ChunkSize = from.ChunkSize
	}

	if len(up.ResumeApiVersion) == 0 {
		up.ResumeApiVersion = from.ResumeApiVersion
	}

	if up.FileConcurrentParts == 0 {
		up.FileConcurrentParts = from.FileConcurrentParts
	}

	up.Tasks.merge(&from.Tasks)
	up.Retry.merge(&from.Retry)
}

func (d *Download) merge(from *Download) {
	if from == nil {
		return
	}

	d.Tasks.merge(&from.Tasks)
	d.Retry.merge(&from.Retry)
}

func (t *Tasks) merge(from *Tasks) {
	if from == nil {
		return
	}

	if t.ConcurrentCount == 0 {
		t.ConcurrentCount = from.ConcurrentCount
	}

	if len(t.StopWhenOneTaskFailed) == 0 {
		t.StopWhenOneTaskFailed = from.StopWhenOneTaskFailed
	}
}

func (r *Retry) merge(from *Retry) {
	if from == nil {
		return
	}

	if r.Max == 0 {
		r.Max = from.Max
	}

	if r.Interval == 0 {
		r.Interval = from.Interval
	}
}
