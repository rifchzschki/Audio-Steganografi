package cli

import (
	"fmt"
	"log"

	"github.com/rifchzschki/Audio-Steganografi/backend/models"
	"github.com/rifchzschki/Audio-Steganografi/backend/service"
)

func Run(args []string) {
	fmt.Println("Running CLI application...")
	switch args[0] {
	case "vigenere":
		ExtendedVignereExample()
	case "stegano":
		SteganoWithLSBExample()
	default:
		fmt.Println("Unknown command")
	}
}

func ExtendedVignereExample() {
	fmt.Println("Encrypting using Extended Vigenere Cipher...")

	key := "mysecret"
	plainText := "Hello, Vigenere 123!"

	cipher := service.NewExtendedVigenereCipher(key)

	encrypted := cipher.Encrypt([]byte(plainText))
	decrypted := cipher.Decrypt(encrypted)

	fmt.Println("Plaintext :", plainText)
	fmt.Println("Encrypted :", encrypted)      
	fmt.Println("Encrypted (string):", string(encrypted)) 
	fmt.Println("Decrypted :", string(decrypted))
}

func SteganoWithLSBExample() {
	fmt.Println("Embedding using Stegano with LSB...")
	
	config := models.LSBConfig{
		Key:            "my-secret-key",
		UseEncryption:  false,
		UseRandomStart: true,
		LSBBits:        1,
		SecretFilename: "secret.txt",
	}

	lsb := service.NewSteganoWithLSB(config)

	cover := make([]byte, 10000) 
	for i := 0; i < len(cover); i++ {
		cover[i] = byte(i % 256)
	}

	secret := []byte("Hello, Steganography ashoy!")

	stegoData, err := lsb.Embed(cover, secret)
	if err != nil {
		log.Fatalf("Embed error: %v", err)
	}
	fmt.Println("Embed berhasil, panjang cover:", len(cover), "panjang stego:", len(stegoData))
	// fmt.Println("Embed berhasil, cover:", string(cover), "stego:", string(stegoData))

	extracted, err := lsb.Extract(stegoData, len(secret))
	if err != nil {
		log.Fatalf("Extract error: %v", err)
	}
	fmt.Println("Extracted data:", string(extracted))
}