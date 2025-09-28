package service

type ExtendedVigenereCipher struct {
	Key []byte
}

func NewExtendedVigenereCipher(key string) *ExtendedVigenereCipher {
	return &ExtendedVigenereCipher{Key: []byte(key)}
}

func (e *ExtendedVigenereCipher) Encrypt(plaintext []byte) []byte {
	ciphertext := make([]byte, len(plaintext))
	keyLen := len(e.Key)
	for i, p := range plaintext {
		k := e.Key[i%keyLen]
		ciphertext[i] = byte((int(p) + int(k)) % 256)
	}
	return ciphertext
}
func (e *ExtendedVigenereCipher) Decrypt(ciphertext []byte) []byte {
	plaintext := make([]byte, len(ciphertext))
	keyLen := len(e.Key)
	for i, ct := range ciphertext {
		k := e.Key[i%keyLen]
		plaintext[i] = byte((int(ct) - int(k) + 256) % 256)
	}
	return plaintext
}