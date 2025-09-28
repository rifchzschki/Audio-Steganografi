package models

type LSBConfig struct {
	Key            string
	UseEncryption  bool
	UseRandomStart bool
	LSBBits        int
	SecretFilename string
}
