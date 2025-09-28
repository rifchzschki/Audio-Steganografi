package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/rifchzschki/Audio-Steganografi/backend/models"
	"github.com/rifchzschki/Audio-Steganografi/backend/utils"
)

type MP3Decoder struct {
	file   *os.File
	data   []byte
	offset int
}

func NewMP3Decoder() *MP3Decoder {
	return &MP3Decoder{}
}

func (d *MP3Decoder) LoadFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	
	d.file = file
	
	data, err := io.ReadAll(file)
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to read file: %v", err)
	}
	
	d.data = data
	d.offset = 0
	
	return nil
}

func (d *MP3Decoder) LoadBytes(data []byte) error {
	if len(data) == 0 {
		return errors.New("data is empty")
	}
	
	d.data = make([]byte, len(data))
	copy(d.data, data)
	d.offset = 0
	
	return nil
}

func (d *MP3Decoder) Close() error {
	if d.file != nil {
		return d.file.Close()
	}
	return nil
}

func (d *MP3Decoder) Decode() (*models.MP3File, error) {
	if len(d.data) == 0 {
		return nil, errors.New("no data loaded")
	}
	
	mp3File := models.NewMP3File()
	d.offset = 0
	
	if err := d.parseID3v2(mp3File); err != nil {
		return nil, fmt.Errorf("failed to parse ID3v2: %v", err)
	}
	
	if err := d.parseFrames(mp3File); err != nil {
		return nil, fmt.Errorf("failed to parse frames: %v", err)
	}
	
	if err := d.parseID3v1(mp3File); err != nil {
		return nil, fmt.Errorf("failed to parse ID3v1: %v", err)
	}
	
	return mp3File, nil
}

func (d *MP3Decoder) parseID3v2(mp3File *models.MP3File) error {
	if len(d.data) < 10 {
		return nil
	}
	
	if !bytes.Equal(d.data[0:3], []byte("ID3")) {
		return nil
	}
	
	header := &models.ID3v2Header{}
	header.Version[0] = d.data[3]
	header.Version[1] = d.data[4]
	header.Flags = d.data[5]
	
	size := 0
	for i := 0; i < 4; i++ {
		size = (size << 7) | int(d.data[6+i]&0x7F)
	}
	header.Size = size
	
	mp3File.ID3v2 = header
	
	if 10+size <= len(d.data) {
		mp3File.ID3v2Data = make([]byte, size)
		copy(mp3File.ID3v2Data, d.data[10:10+size])
	}
	
	d.offset = 10 + size
	return nil
}

func (d *MP3Decoder) parseID3v1(mp3File *models.MP3File) error {
	if len(d.data) < 128 {
		return nil
	}
	
	tagStart := len(d.data) - 128
	if !bytes.Equal(d.data[tagStart:tagStart+3], []byte("TAG")) {
		return nil
	}
	
	tag := &models.ID3v1Tag{}
	tag.Title = string(bytes.TrimRight(d.data[tagStart+3:tagStart+33], "\x00"))
	tag.Artist = string(bytes.TrimRight(d.data[tagStart+33:tagStart+63], "\x00"))
	tag.Album = string(bytes.TrimRight(d.data[tagStart+63:tagStart+93], "\x00"))
	tag.Year = string(bytes.TrimRight(d.data[tagStart+93:tagStart+97], "\x00"))
	tag.Comment = string(bytes.TrimRight(d.data[tagStart+97:tagStart+127], "\x00"))
	tag.Genre = d.data[tagStart+127]
	
	mp3File.ID3v1 = tag
	return nil
}

func (d *MP3Decoder) parseFrames(mp3File *models.MP3File) error {
	maxOffset := len(d.data)
	if mp3File.HasID3v1() {
		maxOffset -= 128
	}
	
	for d.offset < maxOffset-4 {
		if d.offset+4 > len(d.data) {
			break
		}
		
		headerBytes := d.data[d.offset : d.offset+4]
		
		if !d.isValidFrameHeader(headerBytes) {
			d.offset++
			continue
		}
		
		header, err := d.parseFrameHeader(headerBytes)
		if err != nil {
			d.offset++
			continue
		}
		
		if header.FrameLength <= 0 || d.offset+header.FrameLength > maxOffset {
			d.offset++
			continue
		}
		
		frameData := make([]byte, header.FrameLength-4)
		if d.offset+header.FrameLength <= len(d.data) {
			copy(frameData, d.data[d.offset+4:d.offset+header.FrameLength])
		}
		
		frame := &models.MP3Frame{
			Header:      header,
			HeaderBytes: make([]byte, 4),
			Data:        frameData,
		}
		copy(frame.HeaderBytes, headerBytes)
		
		mp3File.AddFrame(frame)
		d.offset += header.FrameLength
	}
	
	return nil
}

func (d *MP3Decoder) isValidFrameHeader(header []byte) bool {
	if len(header) < 4 {
		return false
	}
	
	if header[0] != 0xFF {
		return false
	}
	
	if (header[1] & 0xE0) != 0xE0 {
		return false
	}
	
	versionID := (header[1] >> 3) & 0x03
	if versionID == 1 {
		return false
	}
	
	layer := (header[1] >> 1) & 0x03
	if layer == 0 {
		return false
	}
	
	bitrateIndex := (header[2] >> 4) & 0x0F
	if bitrateIndex == 0 || bitrateIndex == 15 {
		return false
	}
	
	sampleRateIndex := (header[2] >> 2) & 0x03

	return sampleRateIndex != 3
}

func (d *MP3Decoder) parseFrameHeader(headerBytes []byte) (*models.MP3FrameHeader, error) {
	if len(headerBytes) < 4 {
		return nil, errors.New("invalid header length")
	}
	
	header := &models.MP3FrameHeader{}
	
	header.VersionID = int((headerBytes[1] >> 3) & 0x03)
	header.Layer = int((headerBytes[1] >> 1) & 0x03)
	header.ProtectionBit = (headerBytes[1] & 0x01) == 1
	
	bitrateIndex := int((headerBytes[2] >> 4) & 0x0F)
	sampleRateIndex := int((headerBytes[2] >> 2) & 0x03)
	
	header.Padding = (headerBytes[2] & 0x02) != 0
	header.ChannelMode = int((headerBytes[3] >> 6) & 0x03)
	
	versionIndex := 0
	switch header.VersionID {
	case 3:
		versionIndex = 0
	case 2:
		versionIndex = 1
	case 0:
		versionIndex = 2
	}
	
	layerIndex := 0
	
	switch header.Layer {
	case 3:
		layerIndex = 0
	case 2:
		layerIndex = 1
	case 1:
		layerIndex = 2
	}
	
	if versionIndex < 3 && layerIndex < 3 {
		header.Bitrate = utils.MP3BitrateTable[bitrateIndex][layerIndex] * 1000
		header.SampleRate = utils.MP3SampleRateTable[versionIndex][sampleRateIndex]
	}
	
	if header.SampleRate > 0 {
		if header.Layer == 1 {
			header.FrameLength = (12*header.Bitrate/header.SampleRate + int(utils.BoolToInt(header.Padding))) * 4
		} else {
			header.FrameLength = 144*header.Bitrate/header.SampleRate + int(utils.BoolToInt(header.Padding))
		}
	}
	
	return header, nil
}



func (d *MP3Decoder) ExtractAudioData(mp3File *models.MP3File) []byte {
	var audioData []byte
	
	for _, frame := range mp3File.Frames {
		audioData = append(audioData, frame.GetData()...)
	}
	
	return audioData
}

func (d *MP3Decoder) GetAudioDataForSteganography(mp3File *models.MP3File) []byte {
	audioData := d.ExtractAudioData(mp3File)
	
	if len(audioData) == 0 {
		return nil
	}
	
	filteredData := make([]byte, 0, len(audioData))
	for _, b := range audioData {
		if b != 0x00 && b != 0xFF {
			filteredData = append(filteredData, b)
		}
	}
	
	return filteredData
}

func (d *MP3Decoder) GetMetadata(mp3File *models.MP3File) *models.AudioMetadata {
	if len(mp3File.Frames) == 0 {
		return &models.AudioMetadata{}
	}
	
	firstFrame := mp3File.Frames[0]
	header := firstFrame.GetHeader()
	
	channels := 2
	if header.ChannelMode == 3 {
		channels = 1
	}
	
	totalBytes := 0
	for _, frame := range mp3File.Frames {
		totalBytes += len(frame.GetData())
	}
	
	duration := 0.0
	if header.SampleRate > 0 && header.Bitrate > 0 {
		duration = float64(totalBytes*8) / float64(header.Bitrate)
	}
	
	return &models.AudioMetadata{
		SampleRate: header.SampleRate,
		Channels:   channels,
		BitDepth:   16,
		Duration:   duration,
		TotalBytes: totalBytes,
	}
}

func (d *MP3Decoder) EmbedDataWithSteganography(mp3File *models.MP3File, secretData []byte, config models.LSBConfig) *models.SteganographyResult {
	result := &models.SteganographyResult{
		OriginalSize: len(d.ExtractAudioData(mp3File)),
	}
	
	audioData := d.GetAudioDataForSteganography(mp3File)
	if len(audioData) == 0 {
		result.Error = errors.New("no suitable audio data for steganography")
		return result
	}
	
	stego := NewSteganoWithLSB(config)
	if err := stego.ValidateConfig(); err != nil {
		result.Error = fmt.Errorf("invalid config: %v", err)
		return result
	}
	
	capacity := stego.GetCapacity(len(audioData))
	if len(secretData) > capacity {
		result.Error = fmt.Errorf("secret data too large: need %d bytes, capacity %d bytes", len(secretData), capacity)
		return result
	}
	
	modifiedAudioData, err := stego.Embed(audioData, secretData)
	if err != nil {
		result.Error = fmt.Errorf("steganography embedding failed: %v", err)
		return result
	}
	
	modifiedMP3File := d.reconstructMP3File(mp3File, audioData, modifiedAudioData)
	
	result.ModifiedMP3File = modifiedMP3File
	result.ModifiedSize = len(d.ExtractAudioData(modifiedMP3File))
	result.Success = true
	
	return result
}

func (d *MP3Decoder) ExtractDataWithSteganography(mp3File *models.MP3File, dataLength int, config models.LSBConfig) ([]byte, error) {
	audioData := d.GetAudioDataForSteganography(mp3File)
	if len(audioData) == 0 {
		return nil, errors.New("no suitable audio data for steganography")
	}
	
	stego := NewSteganoWithLSB(config)
	if err := stego.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("invalid config: %v", err)
	}
	
	extractedData, err := stego.Extract(audioData, dataLength)
	if err != nil {
		return nil, fmt.Errorf("steganography extraction failed: %v", err)
	}
	
	return extractedData, nil
}

func (d *MP3Decoder) reconstructMP3File(originalMP3 *models.MP3File, originalAudio, modifiedAudio []byte) *models.MP3File {
	newMP3 := models.NewMP3File()
	
	newMP3.ID3v2 = originalMP3.ID3v2
	if originalMP3.ID3v2Data != nil {
		newMP3.ID3v2Data = make([]byte, len(originalMP3.ID3v2Data))
		copy(newMP3.ID3v2Data, originalMP3.ID3v2Data)
	}
	newMP3.ID3v1 = originalMP3.ID3v1
	
	modifiedIndex := 0
	
	for _, originalFrame := range originalMP3.Frames {
		originalData := originalFrame.GetData()
		newData := make([]byte, len(originalData))
		copy(newData, originalData)
		
		for i, b := range originalData {
			if b != 0x00 && b != 0xFF && modifiedIndex < len(modifiedAudio) {
				newData[i] = modifiedAudio[modifiedIndex]
				modifiedIndex++
			}
		}
		
		newFrame := &models.MP3Frame{
			Header:      originalFrame.GetHeader(),
			HeaderBytes: make([]byte, len(originalFrame.(*models.MP3Frame).HeaderBytes)),
			Data:        newData,
		}
		copy(newFrame.HeaderBytes, originalFrame.(*models.MP3Frame).HeaderBytes)
		
		newMP3.AddFrame(newFrame)
	}
	
	return newMP3
}



func (d *MP3Decoder) GetRawData() []byte {
	return d.data
}

func (d *MP3Decoder) GetOffset() int {
	return d.offset
}


func (d *MP3Decoder) GetFrameData(mp3File *models.MP3File, frameIndex int) []byte {
	if frameIndex < 0 || frameIndex >= len(mp3File.Frames) {
		return nil
	}
	
	return mp3File.Frames[frameIndex].GetData()
}

func (d *MP3Decoder) SetFrameData(mp3File *models.MP3File, frameIndex int, data []byte) error {
	if frameIndex < 0 || frameIndex >= len(mp3File.Frames) {
		return errors.New("frame index out of range")
	}
	
	frame, ok := mp3File.Frames[frameIndex].(*models.MP3Frame)
	if !ok {
		return errors.New("invalid frame type")
	}
	
	frame.Data = make([]byte, len(data))
	copy(frame.Data, data)
	
	return nil
}

func (d *MP3Decoder) ValidateMP3File(mp3File *models.MP3File) error {
	if mp3File == nil {
		return errors.New("MP3 file is nil")
	}
	
	if len(mp3File.Frames) == 0 {
		return errors.New("no frames found in MP3 file")
	}
	
	for i, frame := range mp3File.Frames {
		if frame == nil {
			return fmt.Errorf("frame %d is nil", i)
		}
		
		if frame.GetHeader() == nil {
			return fmt.Errorf("frame %d has nil header", i)
		}
		
		if len(frame.GetData()) == 0 {
			return fmt.Errorf("frame %d has no data", i)
		}
	}
	
	return nil
}

func (d *MP3Decoder) GetCompatibleFramesForSteganography(mp3File *models.MP3File) []int {
	var compatibleFrames []int
	
	for i, frame := range mp3File.Frames {
		data := frame.GetData()
		if len(data) > 32 {
			suitableBytes := 0
			for _, b := range data {
				if b != 0x00 && b != 0xFF {
					suitableBytes++
				}
			}
			
			if float64(suitableBytes)/float64(len(data)) > 0.5 {
				compatibleFrames = append(compatibleFrames, i)
			}
		}
	}
	
	return compatibleFrames
}
