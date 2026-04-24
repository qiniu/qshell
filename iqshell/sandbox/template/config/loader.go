package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Load 从给定路径加载并解析 TOML 配置文件。
// 文件不存在时返回 (nil, nil)；解析失败时返回错误。
func Load(path string) (*FileConfig, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve config path: %w", err)
	}

	data, err := os.ReadFile(abs)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read config %s: %w", abs, err)
	}

	cfg := &FileConfig{}
	md, err := toml.NewDecoder(bytes.NewReader(data)).Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("parse config %s: %w", abs, err)
	}
	cfg.defined = make(map[string]bool, len(md.Keys()))
	for _, k := range md.Keys() {
		cfg.defined[k.String()] = true
	}
	cfg.sourcePath = abs
	return cfg, nil
}

// FindInDir 在指定目录中查找 DefaultFileName。
// 找到返回绝对路径；未找到返回 ("", nil)。
func FindInDir(dir string) (string, error) {
	p := filepath.Join(dir, DefaultFileName)
	info, err := os.Stat(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", fmt.Errorf("stat %s: %w", p, err)
	}
	if info.IsDir() {
		return "", nil
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", fmt.Errorf("resolve %s: %w", p, err)
	}
	return abs, nil
}

// LoadFromCwd 在当前工作目录查找并加载默认配置文件。
// 未找到返回 (nil, nil)，方便调用方判断是否存在。
func LoadFromCwd() (*FileConfig, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getwd: %w", err)
	}
	p, err := FindInDir(cwd)
	if err != nil {
		return nil, err
	}
	if p == "" {
		return nil, nil
	}
	return Load(p)
}
