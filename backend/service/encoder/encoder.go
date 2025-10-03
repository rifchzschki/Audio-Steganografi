package encoder

import (
    "crypto/sha256"
    "encoding/binary"
    "fmt"
    "math/rand"
    "os"
    "path/filepath"

    "github.com/rifchzschki/Audio-Steganografi/backend/service"
    "github.com/rifchzschki/Audio-Steganografi/backend/models/meta"
    "github.com/rifchzschki/Audio-Steganografi/backend/models/mp3"
    "github.com/rifchzschki/Audio-Steganografi/backend/utils/payload"
    "github.com/rifchzschki/Audio-Steganografi/backend/utils/sig"
    "github.com/rifchzschki/Audio-Steganografi/backend/utils/psnr"
)

func seedFromKey(key string) int64 { 
    h := sha256.Sum256([]byte(key))
    return int64(binary.LittleEndian.Uint64(h[:8])) 
}

// EncodeFile embeds a secret file into an MP3 file using steganography
func EncodeFile(inputMP3, secretFile, outputMP3, key string, width int, encrypt, random bool) (outputfile string, psnrVal float64,audioQuality string,err error) {
    // Validate width parameter
    if width != 1 && width != 2 && width != 4 && width != 3  {
        return "",0.0,"",fmt.Errorf("width must be 1, 2, 3, or 4")
    }

    // Read cover MP3 file
    coverBytes, err := os.ReadFile(inputMP3)
    if err != nil {
        return "",0.0,"",fmt.Errorf("failed to read input MP3: %v", err)
    }

    // Parse MP3
    f, err := mp3.Parse(coverBytes)
    if err != nil {
        return "",0.0,"",fmt.Errorf("failed to parse MP3: %v", err)
    }

    // Read secret file
    name := filepath.Base(secretFile)
    ext := filepath.Ext(name)
    secretBytes, err := os.ReadFile(secretFile)
    if err != nil {
        return "",0.0,"",fmt.Errorf("failed to read secret file: %v", err)
    }

	fmt.Printf("Encoding file: %s (%s) - %d bytes\n", name, ext, len(secretBytes))

    // Encrypt if requested
    if encrypt {
        secretBytes = service.NewExtendedVigenereCipher(key).Encrypt(secretBytes)
    }

    // Create metadata header
    h := meta.Header{
        Version: 1,
        Flags:   0,
        NLSB:    uint8(width),
        Name:    name,
        Size:    uint64(len(secretBytes)),
        Ext:     ext,
    }
    if encrypt {
        h.Flags |= meta.FlagEncrypted
    }
    if random {
        h.Flags |= meta.FlagRandomStart
    }

    // Pack metadata
    metaBytes := meta.Pack(h)

    // Create bit stream with signature, width, metadata, and payload
    S := sig.Map[width]
    bits := make([]uint8, 0)
    bits = append(bits, S.S...)
    bits = append(bits, sig.WidthByte(width)...)
    bits = append(bits, payload.ToBits(metaBytes)...)
    bits = append(bits, payload.ToBits(secretBytes)...)
    bits = append(bits, S.E...)

    // Extract audio data
    var audio []byte
    for _, fr := range f.Frames {
        audio = append(audio, fr.Data...)
    }

	originalAudio := make([]byte, len(audio))
    copy(originalAudio, audio)

    // Create order array
    order := make([]int, len(audio))
    for i := range audio {
        order[i] = i
    }
    if len(order) == 0 {
        return "",0.0,"",fmt.Errorf("no audio bytes found")
    }

    // Check capacity
    capBits := len(order) * width
    if capBits < len(bits) {
        return "",0.0,"",fmt.Errorf("capacity too small: need %d bits, have %d", len(bits), capBits)
    }

    // Randomize order if requested
    if random {
        rsrc := rand.New(rand.NewSource(seedFromKey(key)))
        for i := len(order) - 1; i > 0; i-- {
            j := rsrc.Intn(i + 1)
            order[i], order[j] = order[j], order[i]
        }
    }

    // Embed bits into audio
    mask := byte((1 << uint(width)) - 1)
    bi := 0
    need := (len(bits) + width - 1) / width

    for t := 0; t < need; t++ {
        pos := order[t]
        var pv byte = 0
        for i := 0; i < width && bi < len(bits); i++ {
            pv = (pv << 1) | bits[bi]
            bi++
        }
        if bi%width != 0 {
            pv <<= uint(width - (bi % width))
        }
        audio[pos] = (audio[pos] &^ mask) | (pv & mask)
    }

    // Update frames with modified audio data
    idx := 0
    for _, fr := range f.Frames {
        for k := range fr.Data {
            if idx < len(audio) {
                fr.Data[k] = audio[idx]
                idx++
            }
        }
    }

    // Serialize and write output
    outBytes := mp3.Serialize(f)
    if err := os.WriteFile(outputMP3, outBytes, 0644); err != nil {
        return "",0.0,"",fmt.Errorf("failed to write output file: %v", err)
    }

    fmt.Printf("Successfully encoded: bits=%d width=%d file=%s\n", len(bits), width, outputMP3)
	psnrValue, _, err := psnr.DetectAudioFormat(originalAudio, audio)
    if err != nil {
        fmt.Printf("Warning: Failed to calculate PSNR: %v\n", err)
        return outputMP3, 0.0, "Unknown", nil
    } else {
        qualityStatus := psnr.GetQualityStatus(psnrValue)
        fmt.Printf("\nQuality Status: %s", qualityStatus)
        return outputMP3, psnrValue, qualityStatus, nil
    }
}
