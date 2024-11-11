package text

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var env = map[string]string{
	"foo": "bar",
}

func getEnv(key string) string {
	var k, def = key, ""
	if n := strings.Index(key, ":"); n > 0 {
		k = key[:n]
		def = key[n+1:]
	}
	if v := env[k]; v != "" {
		return v
	}
	return def
}

func TestExpandEnv(t *testing.T) {
	var expected = map[string]string{
		"aaa$":           "aaa$",
		"$aaa":           "$aaa",
		"aaa${":          "aaa${",
		"${aaa":          "${aaa",
		"$":              "$",
		"$$":             "$$",
		"${foo}":         "bar",
		"${{}":           "",
		"${}":            "",
		"${key}":         "",
		"${key:val}":     "val",
		"${key:val:val}": "val:val",
		"${${}":          "${",
		"${$${}":         "${$",
		"${${}}":         "${}",
		"${a${}a}":       "${aa}",
		"${aa${}}":       "${aa}",
		"${${}aa}":       "${aa}",
		"$${}{aa}":       "${aa}",
	}

	for k, v := range expected {
		assert.Equal(t, v, Expand(k, getEnv))
	}
}
