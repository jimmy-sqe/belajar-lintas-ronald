//go:build ignore

// Helper script: generate bcrypt hashes for seed files.
//
// Run with:
//
//	go run scripts/gen-seed-hash.go
//
// Output is printed to stdout; paste into the seed SQL/JSON files.
package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	passwords := []struct {
		Email    string
		Password string
	}{
		{"demo1@example.com", "Demo1234!"},
		{"demo2@example.com", "Demo1234!"},
	}

	for _, p := range passwords {
		hash, err := bcrypt.GenerateFromPassword([]byte(p.Password), 12)
		if err != nil {
			fmt.Printf("ERR %s: %v\n", p.Email, err)
			continue
		}
		fmt.Printf("%s -> %s\n", p.Email, string(hash))
	}
}
