package values

type Identifier interface {
	// IdentityIc40tk get the id, 'Ic40tk' is used to avoid name conflict
	IdentityIc40tk() int64
}

type StringIdentifier interface {
	// IdentityLMhIA1 get the id, 'LMhIA1' is used to avoid name conflict
	IdentityLMhIA1() string
}

func Identities[T Identifier](v ...T) []int64 {
	ids := make([]int64, len(v))
	for i := range ids {
		ids[i] = v[i].IdentityIc40tk()
	}
	return ids
}

func StringIdentities[T StringIdentifier](v ...T) []string {
	ids := make([]string, len(v))
	for i := range ids {
		ids[i] = v[i].IdentityLMhIA1()
	}
	return ids
}
