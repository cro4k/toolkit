package configuration

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

type localDriver struct {
	root string
}

func NewLocalDriver(root string) Driver {
	if root == "" {
		root = "."
	}
	return &localDriver{root: strings.TrimRight(root, "/") + "/"}
}

func (d *localDriver) Load(_ context.Context, key string) ([]byte, string, error) {
	filename := d.root + strings.TrimLeft(key, ".")
	ext := strings.TrimLeft(filepath.Ext(filename), ".")
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, "", err
	}
	return data, ext, err
}

type localDriverBuilder struct{}

func (b *localDriverBuilder) Build(script string) (Driver, error) {
	return NewLocalDriver(script), nil
}

func NewLocalDriverBuilder() DriverBuilder {
	return &localDriverBuilder{}
}
