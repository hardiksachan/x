// Package xtest provides helper functions for testing.
package xtest

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet   = "abcdefghijklmnopqrstuvwxyz"
	hoursInDay = 24
)

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomString6 generates a random string of length 6
func RandomString6() string {
	const n = 6
	return RandomString(n)
}

// RandomEmailString generates a random email
func RandomEmailString() string {
	return fmt.Sprintf("%s@email.com", RandomString6())
}

// RandomPasswordString generates a random password
func RandomPasswordString() string {
	return fmt.Sprintf("%sA1$", RandomString6())
}

// RandomURL generates a random URL
func RandomURL() string {
	return fmt.Sprintf("https://%s.com", RandomString6())
}

// RandomDate generates a random date
func RandomDate() time.Time {
	startDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, time.December, 31, 0, 0, 0, 0, time.UTC)

	days := int(endDate.Sub(startDate).Hours() / hoursInDay)

	// Generate a random number of days offset from the start date
	offset := rand.Intn(days)

	// Create and return the random date
	randomDate := startDate.Add(time.Duration(offset) * hoursInDay * time.Hour)
	return randomDate
}

// RandomStringArray generates a random string array
func RandomStringArray(size, stringLength int) []string {
	var result []string
	for i := 0; i < size; i++ {
		result = append(result, RandomString(stringLength))
	}
	return result
}
