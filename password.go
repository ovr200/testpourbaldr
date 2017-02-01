package main

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/bcrypt"
)

var chars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!")

// Hash les pass
func encrypt(pass string) string {
	password := []byte(pass)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	result := string(hashedPassword)
	return result
}

//Compare le mot de passe en clair avec le hash , retourne true ou false
func passwordverif(hashpass, pass string) bool {
	// Comparing the password with the hash
	err := bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(pass))
	if err == nil {
		return true
	} else {
		return false
	}
}

func generatepass() string {
	length := 10
	new_pword := make([]byte, length)
	random_data := make([]byte, length+(length/4))
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, random_data); err != nil {
			panic(err)
		}
		for _, c := range random_data {
			if c >= maxrb {
				continue
			}
			new_pword[i] = chars[c%clen]
			i++
			if i == length {
				return string(new_pword)
			}
		}
	}
	panic("unreachable")
}
