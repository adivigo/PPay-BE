package lib

import (
	"fmt"
	"os"

	"github.com/pilinux/argon2"
)

func VerifyHash(hash string, password string) bool {
	result, err := argon2.ComparePasswordAndHash(hash, os.Getenv("HASHKEY"), password)
	fmt.Println("Error during verification:", err)
	return result
}

func GenerateHash(password string) string {
	result, _ := argon2.CreateHash(password, os.Getenv("HASHKEY"), argon2.DefaultParams)
	return result
}
