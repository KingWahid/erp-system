package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	passwords := []string{"Admin@123", "Password@123"}
	for _, p := range passwords {
		h, _ := bcrypt.GenerateFromPassword([]byte(p), 10)
		fmt.Printf("%-15s → %s\n", p, string(h))
	}
}
