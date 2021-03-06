package utils

import (
	"fmt"
	"time"
	"errors"
	"strings"
	"math/rand"
	"encoding/hex"
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
)


func _gen_salt(length int) string{
	BASE_STR := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	salt := ""
	rand.Seed(time.Now().Unix())
	for i:=0; i<length; i++ {
		salt += string(BASE_STR[rand.Intn(len(BASE_STR))])
	}
	return salt
}


func _parse(data string) (method string, salt string, h string){
	r := strings.Split(data, "$")
	if(len(r)) < 3{
		return "", "", ""
	}
	return r[0], r[1], r[2]
}


func _hash_internal(password string, salt string, iter int) (string, error){
	t := pbkdf2.Key([]byte(password), []byte(salt), iter, 32, sha256.New)
	return fmt.Sprintf("pbkdf2:sha256:150000$%s$%s", salt, hex.EncodeToString(t)), nil
}


func GeneratePasswordHash(password string) (string, error){
	salt := _gen_salt(8)
	if len(salt) <= 0 {
		return "", errors.New("gen salt error")
	}
	return _hash_internal(password, salt, 150000)
}


func CheckPasswordHash(pwhash string, password string) bool{
	_,salt,_ := _parse(pwhash)
	t, err := _hash_internal(password, salt, 150000)
	if err != nil {
		return false
	}
	return strings.EqualFold(t, pwhash)
}