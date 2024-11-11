package configuration

import (
	"context"
	"encoding/json"

	"gopkg.in/yaml.v3"
)

const (
	YAML    = "yaml"
	JSON    = "json"
	UNKNOWN = ""

	LocalDriverName = "local"
	RedisDriverName = "redis"
)

var (
	drivers = map[string]Driver{
		LocalDriverName: NewLocalDriver("."),
	}

	builders = map[string]DriverBuilder{}

	unmarshalers = map[string]Unmarshaler{
		UNKNOWN: json.Unmarshal,
		YAML:    yaml.Unmarshal,
		JSON:    json.Unmarshal,
	}
)

type (
	Driver interface {
		Load(ctx context.Context, key string) ([]byte, string, error)
	}

	DriverBuilder interface {
		Build(script string) (Driver, error)
	}

	Unmarshaler func([]byte, any) error
)

func SetDriver(name string, d Driver) {
	drivers[name] = d
}

func SetDriverBuilder(name string, d DriverBuilder) {
	builders[name] = d
}

func SetUnmarshaler(contentType string, unmarshaler Unmarshaler) {
	unmarshalers[contentType] = unmarshaler
}
