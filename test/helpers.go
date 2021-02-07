package test

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

var tlds = [5]string{"com", "me", "net", "gov", "edu"}

// Generates a random domain name.
func RandomDomainName(length int) string {
	rand.Seed(time.Now().UnixNano())
	charSet := "abcdefghijklmnopqrstuvwxyz"
	var output strings.Builder
	defer output.Reset()

	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}

	random := rand.Intn(len(tlds))
	randomTld := tlds[random]

	output.WriteString(".")
	output.WriteString(randomTld)

	return output.String()
}

func CheckError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
