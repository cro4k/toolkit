package configuration

import "github.com/cro4k/toolkit/text"

func ExpandEnvUnmarshaler(unmarshaler Unmarshaler) Unmarshaler {
	return func(data []byte, dst any) error {
		data = []byte(text.ExpandEnv(string(data)))
		return unmarshaler(data, dst)
	}
}

// EnableExpandEnv make it can access environment variable in config content.
// For example:
//
//		redis:
//		  host: ${REDIS_HOST}
//	      port: ${REDIS_PORT:6379}
func EnableExpandEnv() {
	for k, v := range unmarshalers {
		unmarshalers[k] = ExpandEnvUnmarshaler(v)
	}
}
