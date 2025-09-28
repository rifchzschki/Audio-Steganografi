package cli

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"os"
	"github.com/rifchzschki/Audio-Steganografi/backend/models"
	"github.com/rifchzschki/Audio-Steganografi/backend/service"
	"github.com/rifchzschki/Audio-Steganografi/backend/service/decoder"
	"github.com/rifchzschki/Audio-Steganografi/backend/service/encoder"
)

func Run(args []string) {
	fmt.Println("Running CLI application...")
	switch args[0] {
	case "vigenere":
		ExtendedVignereExample()
	case "stegano":
		SteganoWithLSBExample()
	case "stegano-audio":
		SteganoAudioExample()
	case "en-x":
		EncodeX()
	case "dec-x":
		DecodeX()
	default:
		fmt.Println("Unknown command")
	}
}

func DecodeX(){
	// Example usage with manual variables
    inputFile := "stego.mp3"     // Input steganographic MP3 file
    outputFile := "decoded.txt"  // Output file (can be empty to use original name)
    key := "STEGANO"            // Encryption key/seed
    random := true              // Use random order
    debug := false              // Enable debug logs
    
    if err := decoder.DecodeFile(inputFile, outputFile, key, random, debug); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}

func EncodeX(){
	// Example usage with manual variables
    inputMP3 := "cover.mp3"      // Input MP3 cover file
    secretFile := "secret.txt"   // Secret file to hide
    outputMP3 := "stego.mp3"     // Output steganographic MP3
    key := "STEGANO"             // Encryption key/seed
    width := 1                   // LSB width (1, 2, or 4)
    encrypt := false             // Encrypt the payload
    random := true               // Use random order

    if err := encoder.EncodeFile(inputMP3, secretFile, outputMP3, key, width, encrypt, random); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
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

func SteganoAudioExample() {
	fmt.Println("=== Real MP3 Steganography Test ===")
	
	samplePath := filepath.Join("sample", "sample-6s.mp3")
	
	secretData := []byte("Hellowwworldkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkjjjj")
	config := models.LSBConfig{
		LSBBits: 2,
		Key:     "test",
	}
	
	fmt.Printf("Original secret: %v\n", secretData)
	fmt.Printf("Original secret as string: %q\n", secretData)
	
	decoder := service.NewMP3Decoder()
	defer decoder.Close()
	
	err := decoder.LoadFile(samplePath)
	if err != nil {
		log.Fatal("LoadFile failed:", err)
	}
	
	mp3File, err := decoder.Decode()
	if err != nil {
		log.Fatal("Decode failed:", err)
	}
	
	fmt.Printf("MP3 decoded successfully, frames: %d\n", mp3File.GetTotalFrames())
	
	originalAudioData := decoder.GetAudioDataForSteganography(mp3File)
	fmt.Printf("Original audio data length: %d\n", len(originalAudioData))
	
	stego := service.NewSteganoWithLSB(config)
	capacity := stego.GetCapacity(len(originalAudioData))
	fmt.Printf("Steganography capacity: %d bytes\n", capacity)
	
	if len(secretData) > capacity {
		log.Fatal("Secret data too large for real MP3")
	}
	
	result := decoder.EmbedDataWithSteganography(mp3File, secretData, config)
	if !result.Success {
		log.Fatal("MP3 embedding failed:", result.Error)
	}
	
	fmt.Printf("MP3 embedding successful\n")
	
	extractedFromMP3, err := decoder.ExtractDataWithSteganography(result.ModifiedMP3File, len(secretData), config)
	if err != nil {
		log.Fatal("MP3 extraction failed:", err)
	}
	
	fmt.Printf("Extracted from MP3: %v\n", extractedFromMP3)
	fmt.Printf("Extracted from MP3 as string: %q\n", extractedFromMP3)
	fmt.Printf("Real MP3 steganography works: %v\n", bytes.Equal(secretData, extractedFromMP3))
	
	if bytes.Equal(secretData, extractedFromMP3) {
		fmt.Println("SUCCESS: Real MP3 steganography is working perfectly!")
	} else {
		fmt.Println("FAILED: Real MP3 steganography has issues")
		
		fmt.Println("\nTrying with smaller data...")
		smallSecret := []byte("Hi")
		smallResult := decoder.EmbedDataWithSteganography(mp3File, smallSecret, config)
		if !smallResult.Success {
			log.Fatal("Small embedding failed:", smallResult.Error)
		}
		
		smallExtracted, err := decoder.ExtractDataWithSteganography(smallResult.ModifiedMP3File, len(smallSecret), config)
		if err != nil {
			log.Fatal("Small extraction failed:", err)
		}
		
		fmt.Printf("Small test - Original: %q, Extracted: %q, Success: %v\n", 
			smallSecret, smallExtracted, bytes.Equal(smallSecret, smallExtracted))
	}
}