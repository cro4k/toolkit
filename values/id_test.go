package values

import (
	"strconv"
	"testing"
)

type Item struct {
	id int64
}

func (i *Item) IdentityIc40tk() int64  { return i.id }
func (i *Item) IdentityLMhIA1() string { return strconv.FormatInt(i.id, 10) }

func TestIdentifier(t *testing.T) {
	var items = []*Item{{1}, {2}, {3}, {4}, {5}}

	t.Log(Identities(items...))
	t.Log(StringIdentities(items...))
}
