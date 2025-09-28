package models

type LSBConfig struct {
	UseEncryption  bool
	UseRandomStart bool
	LSBBits        int
	SecretFilename string
}

