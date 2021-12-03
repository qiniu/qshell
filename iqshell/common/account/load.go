package account

import "fmt"

var info = &pathInfo{}

type Option func(i *pathInfo)

func AccountPath(path string) Option {
	return func(i *pathInfo) {
		i.accountPath = path
	}
}

func OldAccountPath(path string) Option {
	return func(i *pathInfo) {
		i.oldAccountPath = path
	}
}

func AccountDBPath(path string) Option {
	return func(i *pathInfo) {
		i.accountDBPath = path
	}
}

// 保证 accountPath、oldAccountPath、accountDBPath 均不为空
func Load(options ...Option) error {
	for _, option := range options {
		option(info)
	}

	if info.accountDBPath == "" {
		return fmt.Errorf("empty account db path\n")
	}

	if info.accountPath == "" {
		return fmt.Errorf("empty account path\n")
	}

	if info.oldAccountPath == "" {
		return fmt.Errorf("empty old account db path\n")
	}

	return nil
}

type pathInfo struct {
	accountPath    string
	oldAccountPath string
	accountDBPath  string
}
