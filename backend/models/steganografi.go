package models

type LSBConfig struct {
	AudioFileName  string
	SecretFilename string
	UseEncryption  bool
	UseRandomStart bool
	LSBBits        int
	Key            string
}

type LSBConfigDecoder struct {
	Key            string 
	UseRandomStart bool
	SecretFilename string
	OutputFileName string
}

type SteganographyResult struct {
	ModifiedMP3File *MP3File
	OriginalSize    int
	ModifiedSize    int
	Success         bool
	Error           error
}
