package util

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashSaltPassword(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(hash), nil
}

func ComparePasswords(hashedPwd string, password []byte) bool {
	hash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
