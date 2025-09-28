package models

type LSBConfig struct {
	Key            string
	UseEncryption  bool
	UseRandomStart bool
	LSBBits        int
	SecretFilename string
}

type SteganographyResult struct {
	ModifiedMP3File *MP3File
	OriginalSize    int
	ModifiedSize    int
	Success         bool
	Error           error
}
