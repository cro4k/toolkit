package values

import (
	"fmt"
	"strconv"
)

type Integer interface {
	int | int8 | int16 | int32 | int64
}

type UInteger interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type Float interface {
	float32 | float64
}

type Number interface {
	Integer | UInteger | Float
}

func ParseInteger[T Integer](s string) T {
	val, _ := strconv.ParseInt(s, 10, 64)
	return T(val)
}

func ParseUInteger[T UInteger](s string) T {
	val, _ := strconv.ParseUint(s, 10, 64)
	return T(val)
}

func ParseFloat[T Float](s string) T {
	val, _ := strconv.ParseFloat(s, 64)
	return T(val)
}

func NumbersTo[A, B Number](from []A) []B {
	results := make([]B, 0, len(from))
	for _, a := range from {
		results = append(results, B(a))
	}
	return results
}

func Iterate[T any](list []T, f func(v T)) {
	for _, v := range list {
		f(v)
	}
}

func IterateX[A, B any](list []A, f func(v A) B) []B {
	results := make([]B, 0, len(list))
	for _, a := range list {
		results = append(results, f(a))
	}
	return results
}

func Strings[T fmt.Stringer](list []T) []string {
	return IterateX(list, func(v T) string { return v.String() })
}
