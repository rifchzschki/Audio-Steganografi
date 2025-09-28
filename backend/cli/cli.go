package cli

import (
	"fmt"

	"github.com/rifchzschki/Audio-Steganografi/backend/service"
)

func Run(args []string) {
	fmt.Println("Running CLI application...")
	switch args[0] {
	case "vigenere":
		ExtendedVignereExample()
	default:
		fmt.Println("Unknown command")
	}
}

func ExtendedVignereExample() {
	fmt.Println("Encrypting using Extended Vigenere Cipher...")
	// Implementasi enkripsi di sini

	key := "mysecret"
	plainText := "Hello, Vigenere 123!"

	cipher := service.NewExtendedVigenereCipher(key)

	encrypted := cipher.Encrypt([]byte(plainText))
	decrypted := cipher.Decrypt(encrypted)

	fmt.Println("Plaintext :", plainText)
	fmt.Println("Encrypted :", encrypted)       // dalam bentuk byte array
	fmt.Println("Encrypted (string):", string(encrypted)) // bisa juga dipaksa jadi string
	fmt.Println("Decrypted :", string(decrypted))
}