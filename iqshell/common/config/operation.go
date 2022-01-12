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

	c.Hosts.merge(from.Hosts)
	c.Up.merge(from.Up)
	c.Download.merge(from.Download)
}

func (h *Hosts) merge(from *Hosts) {
	if from == nil {
		return
	}

	if len(h.UC) == 0 {
		h.UC = from.UC
	}

	if len(h.Api) == 0 {
		h.Api = from.Api
	}

	if len(h.Rsf) == 0 {
		h.Rsf = from.Rsf
	}

	if len(h.Rs) == 0 {
		h.Rs = from.Rs
	}

	if len(h.Io) == 0 {
		h.Io = from.Io
	}

	if len(h.Up) == 0 {
		h.Up = from.Up
	}
}

func (up *Up) merge(from *Up) {
	if from == nil {
		return
	}

	if up.PutThreshold == 0 {
		up.PutThreshold = from.PutThreshold
	}

	up.Tasks.merge(from.Tasks)
	up.Retry.merge(from.Retry)
}

func (d *Download) merge(from *Download) {
	if from == nil {
		return
	}

	d.Tasks.merge(from.Tasks)
	d.Retry.merge(from.Retry)
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
