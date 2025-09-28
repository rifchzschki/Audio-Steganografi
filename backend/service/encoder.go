package service

import (
	"fmt"
	"os"

	"github.com/rifchzschki/Audio-Steganografi/backend/models"
)

type MP3Encoder struct{}

func NewMP3Encoder() *MP3Encoder {
	return &MP3Encoder{}
}

func (e *MP3Encoder) EncodeToBytes(mp3File *models.MP3File) ([]byte, error) {
	var result []byte

	if mp3File.HasID3v2() && mp3File.ID3v2Data != nil {
		id3v2Header := make([]byte, 10)
		copy(id3v2Header[0:3], []byte("ID3"))
		id3v2Header[3] = mp3File.ID3v2.Version[0]
		id3v2Header[4] = mp3File.ID3v2.Version[1]
		id3v2Header[5] = mp3File.ID3v2.Flags

		size := mp3File.ID3v2.Size
		for i := 3; i >= 0; i-- {
			id3v2Header[6+(3-i)] = byte(size >> (i * 7) & 0x7F)
		}

		result = append(result, id3v2Header...)
		result = append(result, mp3File.ID3v2Data...)
	}

	for _, frame := range mp3File.Frames {
		mp3Frame, ok := frame.(*models.MP3Frame)
		if !ok {
			continue
		}

		result = append(result, mp3Frame.HeaderBytes...)
		result = append(result, mp3Frame.Data...)
	}

	if mp3File.HasID3v1() {
		id3v1Tag := make([]byte, 128)
		copy(id3v1Tag[0:3], []byte("TAG"))

		tag := mp3File.ID3v1
		copy(id3v1Tag[3:33], []byte(tag.GetTitle()))
		copy(id3v1Tag[33:63], []byte(tag.GetArtist()))
		copy(id3v1Tag[63:93], []byte(tag.GetAlbum()))
		copy(id3v1Tag[93:97], []byte(tag.GetYear()))
		copy(id3v1Tag[97:127], []byte(tag.GetComment()))
		id3v1Tag[127] = tag.GetGenre()

		result = append(result, id3v1Tag...)
	}

	return result, nil
}

func (e *MP3Encoder) SaveToFile(mp3File *models.MP3File, filePath string) error {
	data, err := e.EncodeToBytes(mp3File)
	if err != nil {
		return fmt.Errorf("failed to encode MP3: %v", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

