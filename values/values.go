package values

func Default[T any](values []T, def T) (t T) {
	if len(values) > 0 {
		return values[0]
	}
	return def
}

func Then[T any](then bool, a, b T) T {
	if then {
		return a
	}
	return b
}

func ThenSet[T any](then bool, set func() T) T {
	var t T
	if then {
		t = set()
	}
	return t
}

func ThenSetAny[T any](then bool, set any) T {
	var t T
	if then {
		switch val := set.(type) {
		case T:
			t = val
		case func() T:
			t = val()
		}
	}
	return t
}

func ThenFunc[T any](then bool, a, b func() T) T {
	var t T
	if then {
		if a != nil {
			t = a()
		}
	} else {
		if b != nil {
			t = b()
		}
	}
	return t
}

func ThenAny[T any](then, a, b any) T {
	var condition bool
	switch v := then.(type) {
	case bool:
		condition = v
	case func() bool:
		condition = v()
	default:
		condition = false
	}

	var t T
	if condition {
		if val, ok := a.(T); ok {
			return val
		}
		if f, ok := a.(func() T); ok {
			return f()
		}
	} else {
		if val, ok := b.(T); ok {
			return val
		}
		if f, ok := b.(func() T); ok {
			return f()
		}
	}
	return t
}

func Contains[T comparable](list []T, v T) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}
	return false
}
