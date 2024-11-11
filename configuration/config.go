package configuration

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

var (
	configDriver = LocalDriverName
	configKey    = "config.yaml"
	configScript = ""

	ErrUnknownDriver = errors.New("the configuration driver is unknown")
)

func SetFlag() {
	flag.StringVar(&configDriver, "config-driver", LocalDriverName, "Setting config driver")
	flag.StringVar(&configScript, "config-script", "", "Setting config script")
	flag.StringVar(&configKey, "config", "config.yaml", "Setting config key")
}

// Load config
//
// # load by default (local file driver):
// go run main.go
//
// # load from local file driver:
// go run main.go -config-driver local -config config.yaml
//
// # load from pre-registered remote driver:
// go run main.go -config-driver redis -config config.yaml
//
// # load from dynamic created remote driver:
// go run main.go -config-driver consul -config-script consul.yaml -config app/config.yaml
func Load(ctx context.Context, dst any) error {
	if !flag.Parsed() {
		flag.Parse()
	}
	if configScript != "" {
		return LoadWithScript(ctx, configDriver, configScript, configKey, dst)
	}

	driver, ok := drivers[configDriver]
	if !ok {
		return fmt.Errorf("load config failed: driver=%s, key=%s, err=%w", configDriver, configKey, ErrUnknownDriver)
	}
	return LoadWithDriver(ctx, driver, configKey, dst)
}

func LoadWithScript(ctx context.Context, driverName, driverScript, key string, dst any) error {
	builder, ok := builders[driverName]
	if !ok {
		return fmt.Errorf("load config failed: driver=%s, key=%s, err=%w", configDriver, key, ErrUnknownDriver)
	}
	driver, err := builder.Build(driverScript)
	if err != nil {
		return fmt.Errorf("load config failed: driver=%s, key=%s, err=%w", driverName, key, err)
	}
	return LoadWithDriver(ctx, driver, key, dst)
}

func LoadWithDriver(ctx context.Context, driver Driver, key string, dst any) error {
	data, contentType, err := driver.Load(ctx, key)
	if err != nil {
		return fmt.Errorf("load config failed: key=%s, contentType=%s, err=%w",
			key, contentType, err)
	}
	unmarshaler := unmarshalers[contentType]
	if unmarshaler == nil {
		unmarshaler = unmarshalers[UNKNOWN]
	}
	if err = unmarshaler(data, dst); err != nil {
		return fmt.Errorf("load config failed: key=%s, contentType=%s, err=%w",
			key, contentType, err)
	}
	return nil
}
