package randutil

import "math/rand"

const (
	// SeedNumber is the number 0-9
	SeedNumber = "0123456789"
	// SeedLetter is the all lowercase letter
	SeedLetter = "abcdefghijklmnopqrstuvwxyz"
	// SeedCapital is the all uppercase letter
	SeedCapital = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// SeedWords is the all number and lowercase letter.
	SeedWords = SeedNumber + SeedLetter
	// SeedReadFriendly removed the ambiguous words, such as number '0' and letter 'o', and so on.
	SeedReadFriendly = "23456789ABCDEFGHJKMNPQRSTWXYZ"
)

// Text generate a string
func Text(n int, seed ...string) string {
	var s string
	if len(seed) > 0 && seed[0] != "" {
		s = seed[0]
	} else {
		s = SeedWords
	}
	results := make([]byte, n)
	for i := 0; i < n; i++ {
		results[i] = s[rand.Intn(len(s))]
	}
	return string(results)
}

// Bytes generate a bytes slice
func Bytes(n int) []byte {
	results := make([]byte, n)
	for i := 0; i < n; i++ {
		results[i] = byte(rand.Intn(256))
	}
	return results
}

// Hex generate a string that contains only hex number.
func Hex(n int) string {
	return Text(n, SeedNumber+"abcdef")
}
