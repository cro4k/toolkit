package text

import (
	"bytes"
	"os"
	"strings"
)

type expandOptions struct {
	left, right []byte
}

type ExpandOption func(o *expandOptions)

func WithExpandLeft(left []byte) ExpandOption {
	return func(o *expandOptions) {
		o.left = left
	}
}

func WithExpandRight(right []byte) ExpandOption {
	return func(o *expandOptions) {
		o.right = right
	}
}

func applyExpandOptions(opt *expandOptions, opts ...ExpandOption) {
	for _, o := range opts {
		o(opt)
	}
}

func ExpandEnv(text string, options ...ExpandOption) string {
	return Expand(text, GetEnvWithDefault, options...)
}

func Expand(text string, f func(string) string, options ...ExpandOption) string {
	opt := &expandOptions{
		left:  []byte{'$', '{'},
		right: []byte{'}'},
	}
	applyExpandOptions(opt, options...)
	left, right := opt.left, opt.right

	var buf []byte
	var key []byte
	var expanding bool
	var i = 0
	for i < len(text) {
		if len(text) >= i+len(left) && bytes.Equal([]byte(text[i:i+len(left)]), left) {
			if expanding {
				buf = append(buf, left...)
			}
			buf = append(buf, key...)
			key = []byte{}
			expanding = true
			i += len(left)
			continue
		}
		if len(text) >= i+len(right) && expanding && bytes.Equal([]byte(text[i:i+len(right)]), right) {
			buf = append(buf, f(string(key))...)
			key = []byte{}
			expanding = false
			i += len(right)
			continue
		}
		if expanding {
			key = append(key, text[i])
		} else {
			buf = append(buf, text[i])
		}
		i++
	}
	if expanding {
		buf = append(buf, left...)
		buf = append(buf, key...)
	}
	return string(buf)
}

// GetEnvWithDefault get env and return default value if env not found
// eg:
//
//	 // Find the value of 'KEY' in env, and return 'HELLO' if 'KEY' not existed.
//	GetEnvWithDefault("KEY:HELLO")
func GetEnvWithDefault(key string) string {
	var def string
	key, def, _ = strings.Cut(key, ":")
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
