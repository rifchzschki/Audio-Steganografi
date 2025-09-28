package service

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"

	"github.com/rifchzschki/Audio-Steganografi/backend/models"
)
type SteganoWithLSB struct {
	LSBConfig models.LSBConfig
	rng       *rand.Rand
}

func NewSteganoWithLSB(config models.LSBConfig) *SteganoWithLSB {
	seed := deriveSeedFromKey(config.Key)
	source := rand.NewSource(seed)
	return &SteganoWithLSB{
		LSBConfig: config,
		rng:      rand.New(source),
	}
}

func deriveSeedFromKey(key string) int64 {
	hash := sha256.Sum256([]byte(key))
	return int64(binary.LittleEndian.Uint64(hash[:8]))
}

func (lsb *SteganoWithLSB) GetCapacity(coverLen int) int {
	return (coverLen * lsb.LSBConfig.LSBBits) / 8
}

func (lsb *SteganoWithLSB) ValidateConfig() error {
	if lsb.LSBConfig.LSBBits < 1 || lsb.LSBConfig.LSBBits > 4 {
		return errors.New("LSBBits harus 1-4")
	}
	if lsb.LSBConfig.Key == "" {
		return errors.New("key tidak boleh kosong")
	}
	return nil
}


func (lsb *SteganoWithLSB) generatePositions(audioLen, bitsNeeded int) []int {
	if audioLen <= 0 || bitsNeeded <= 0 {
		return []int{}
	}
	
	if bitsNeeded > audioLen {
		bitsNeeded = audioLen
	}
	
	positions := make([]int, 0, bitsNeeded)

	if lsb.LSBConfig.UseRandomStart {
		seed := deriveSeedFromKey(lsb.LSBConfig.Key)
		tempRng := rand.New(rand.NewSource(seed))
		
		used := make(map[int]bool)
		attempts := 0
		maxAttempts := audioLen * 2
		
		for len(positions) < bitsNeeded && attempts < maxAttempts {
			pos := tempRng.Intn(audioLen)
			if !used[pos] {
				positions = append(positions, pos)
				used[pos] = true
			}
			attempts++
		}
		
		if len(positions) < bitsNeeded {
			for i := 0; i < audioLen && len(positions) < bitsNeeded; i++ {
				if !used[i] {
					positions = append(positions, i)
					used[i] = true
				}
			}
		}
	} else {
		for i := 0; i < bitsNeeded && i < audioLen; i++ {
			positions = append(positions, i)
		}
	}

	return positions
}

func (lsb *SteganoWithLSB) Embed(cover []byte, secretData []byte) ([]byte, error) {
	if lsb.LSBConfig.LSBBits < 1 || lsb.LSBConfig.LSBBits > 4 {
		return nil, errors.New("LSBBits harus 1-4")
	}
	
	if len(cover) == 0 {
		return nil, errors.New("cover data kosong")
	}
	
	if len(secretData) == 0 {
		return nil, errors.New("secret data kosong")
	}

	requiredBits := len(secretData) * 8
	capacity := len(cover) * lsb.LSBConfig.LSBBits
	if requiredBits > capacity {
		return nil, fmt.Errorf("cover tidak cukup menampung secret data: butuh %d bits, kapasitas %d bits", requiredBits, capacity)
	}

	stego := make([]byte, len(cover))
	copy(stego, cover)

	bitsNeeded := requiredBits / lsb.LSBConfig.LSBBits
	if requiredBits%lsb.LSBConfig.LSBBits != 0 {
		bitsNeeded++
	}
	positions := lsb.generatePositions(len(cover), bitsNeeded)
	secretBits := make([]byte, requiredBits)
	for i, b := range secretData {
		for j := 0; j < 8; j++ {
			bit := (b >> (7 - j)) & 1
			secretBits[i*8+j] = bit
		}
	}

	mask := byte((1 << lsb.LSBConfig.LSBBits) - 1)
	bitIndex := 0
	
	for _, pos := range positions {
		if bitIndex >= len(secretBits) {
			break
		}
		
		var payloadValue byte = 0
		for i := 0; i < lsb.LSBConfig.LSBBits && bitIndex < len(secretBits); i++ {
			payloadValue = (payloadValue << 1) | secretBits[bitIndex]
			bitIndex++
		}
		
		stego[pos] = (stego[pos] & ^mask) | (payloadValue & mask)
	}

	return stego, nil
}

func (lsb *SteganoWithLSB) Extract(stego []byte, length int) ([]byte, error) {
	if lsb.LSBConfig.LSBBits < 1 || lsb.LSBConfig.LSBBits > 4 {
		return nil, errors.New("LSBBits harus 1-4")
	}
	
	if len(stego) == 0 {
		return nil, errors.New("stego data kosong")
	}
	
	if length <= 0 {
		return nil, errors.New("length harus lebih dari 0")
	}

	requiredBits := length * 8
	bitsNeeded := requiredBits / lsb.LSBConfig.LSBBits
	if requiredBits%lsb.LSBConfig.LSBBits != 0 {
		bitsNeeded++
	}
	
	if bitsNeeded > len(stego) {
		return nil, fmt.Errorf("stego data tidak cukup: butuh %d posisi, tersedia %d", bitsNeeded, len(stego))
	}
	
	positions := lsb.generatePositions(len(stego), bitsNeeded)
	mask := byte((1 << lsb.LSBConfig.LSBBits) - 1)
	extractedBits := make([]byte, 0, requiredBits)
	
	for _, pos := range positions {
		payloadValue := stego[pos] & mask
		
		for i := lsb.LSBConfig.LSBBits - 1; i >= 0; i-- {
			bit := (payloadValue >> i) & 1
			extractedBits = append(extractedBits, bit)
			
			if len(extractedBits) >= requiredBits {
				break
			}
		}
		
		if len(extractedBits) >= requiredBits {
			break
		}
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		var b byte
		for j := 0; j < 8; j++ {
			bitIndex := i*8 + j
			if bitIndex < len(extractedBits) {
				b = (b << 1) | extractedBits[bitIndex]
			}
		}
		result[i] = b
	}

	return result, nil
}