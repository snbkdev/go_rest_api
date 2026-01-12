package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func EncryptString(key, text string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	plaintext := []byte(text)
	fixedIV := []byte("1234567890abcdef")
	cfb := cipher.NewCFBEncrypter(block, fixedIV)
	cipherText := make([]byte, len(plaintext))
	cfb.XORKeyStream(cipherText, plaintext)
	return base64.StdEncoding.EncodeToString(cipherText)
}

func DecryptString(key, text string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	fixedIV := []byte("1234567890abcdef")
	cipherText, _ := base64.StdEncoding.DecodeString(text)
	cfb := cipher.NewCFBDecrypter(block, fixedIV)
	plaintext := make([]byte, len(cipherText))
	cfb.XORKeyStream(plaintext, cipherText)
	return string(plaintext)
}