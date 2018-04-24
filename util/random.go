package util

import (
	"math/rand"
	"time"
	"bytes"
)

// Generate a random int between min and max, inclusive
func Random(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max - min) + min
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

// Returns a unique (ish) id we can attach to resources and tfstate files so they don't conflict with each other
// Uses base 62 to generate a 6 character string that's unlikely to collide with the handful of tests we run in
// parallel. Based on code here: http://stackoverflow.com/a/9543797/483528
func UniqueId() string {

	const BASE_62_CHARS = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	const UNIQUE_ID_LENGTH = 6 // Should be good for 62^6 = 56+ billion combinations

	var out bytes.Buffer

	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < UNIQUE_ID_LENGTH; i++ {
		out.WriteByte(BASE_62_CHARS[generator.Intn(len(BASE_62_CHARS))])
	}

	return out.String()

}