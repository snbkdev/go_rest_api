package main

import (
	"log"
	"restapi/chapter_09/encryptString/utils"
)

func main() {
	key := "111023043350789514532147"
	message := "I am a Message"
	log.Println("Original message: ", message)
	encryptedMessage := utils.EncryptString(key, message)
	log.Println("Encrypted message: ", encryptedMessage)
	decryptedMessage := utils.DecryptString(key, encryptedMessage)
	log.Println("Decrypted message: ", decryptedMessage)
}