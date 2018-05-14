package random

import (
	"bytes"
	"math/rand"
	"time"
)

// Generate a random int between min and max, inclusive
func Random(min int, max int) int {
	return newRand().Intn(max-min) + min
}

// Pick a random element in the slice of ints
func RandomInt(elements []int) int {
	index := Random(0, len(elements))
	return elements[index]
}

// Pick a random element in the slice of string
func RandomString(elements []string) string {
	index := Random(0, len(elements))
	return elements[index]
}

const base62chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const uniqueIdLength = 6 // Should be good for 62^6 = 56+ billion combinations

// Returns a unique (ish) id we can attach to resources and tfstate files so they don't conflict with each other
// Uses base 62 to generate a 6 character string that's unlikely to collide with the handful of tests we run in
// parallel. Based on code here: http://stackoverflow.com/a/9543797/483528
func UniqueId() string {
	var out bytes.Buffer

	generator := newRand()
	for i := 0; i < uniqueIdLength; i++ {
		out.WriteByte(base62chars[generator.Intn(len(base62chars))])
	}

	return out.String()
}

// Create a new random number generator, seeding it with the current system time
func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
