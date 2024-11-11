package cluster

import "strings"

// Metadata
// the key is case-insensitive
type Metadata map[string][]string

func (m Metadata) Add(key string, values ...string) {
	key = strings.ToLower(key)
	m[key] = append(m[key], values...)
}

func (m Metadata) Del(key string) {
	key = strings.ToLower(key)
	delete(m, key)
}
func (m Metadata) Set(key string, value string) {
	key = strings.ToLower(key)
	m[key] = []string{value}
}

func (m Metadata) Get(key string) string {
	key = strings.ToLower(key)
	if val, ok := m[key]; ok {
		return val[0]
	}
	return ""
}

func (m Metadata) Gets(key string) []string {
	key = strings.ToLower(key)
	return m[key]
}
