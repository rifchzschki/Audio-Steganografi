package decoder

import (
    "crypto/sha256"
    "encoding/binary"
    "fmt"
    "math/rand"
    "os"

    "github.com/rifchzschki/Audio-Steganografi/backend/service"
    "github.com/rifchzschki/Audio-Steganografi/backend/models/meta"
    "github.com/rifchzschki/Audio-Steganografi/backend/models/mp3"
    "github.com/rifchzschki/Audio-Steganografi/backend/utils/payload"
    "github.com/rifchzschki/Audio-Steganografi/backend/utils/sig"
)

func seedFromKey(key string) int64 { 
    h := sha256.Sum256([]byte(key))
    return int64(binary.LittleEndian.Uint64(h[:8])) 
}

func find(sigBits, hay []uint8) int { 
    if len(sigBits) == 0 || len(hay) < len(sigBits) { return -1 }
    for i := 0; i+len(sigBits) <= len(hay); i++ { 
        ok := true
        for j := 0; j < len(sigBits); j++ { 
            if hay[i+j] != sigBits[j] { 
                ok = false
                break 
            } 
        }
        if ok { return i } 
    }
    return -1 
}

func tryDecode(audio []byte, key string, random bool, w int, dbg bool) ([]byte, *meta.Header, bool) {
    eligible := make([]int, 0, len(audio))
    for i, b := range audio { 
        if b != 0x00 && b != 0xFF { 
            eligible = append(eligible, i) 
        } 
    }
    if len(eligible) == 0 { return nil, nil, false }
    
    order := make([]int, len(audio))
    for i := range audio { order[i] = i }
    if len(order) == 0 { panic("no audio bytes") }
    
    if random { 
        rsrc := rand.New(rand.NewSource(seedFromKey(key)))
        for i := len(order) - 1; i > 0; i-- { 
            j := rsrc.Intn(i + 1)
            order[i], order[j] = order[j], order[i] 
        } 
    }
    
    mask := byte((1 << uint(w)) - 1)
    stream := make([]uint8, 0, len(order)*w)
    for _, pos := range order { 
        v := audio[pos] & mask
        for i := w - 1; i >= 0; i-- { 
            stream = append(stream, (v>>uint(i))&1) 
        } 
    }
    
    if dbg {
        nshow := 48
        if len(stream) < nshow { nshow = len(stream) }
        fmt.Printf("[DBG] w=%d random=%v firstBits=%v...\n", w, random, stream[:nshow])
        npos := 10
        if len(order) < npos { npos = len(order) }
        fmt.Printf("[DBG] firstPos=%v\n", order[:npos])
        sg := sig.Map[w]
        fmt.Printf("[DBG] sigS=%v\n", sg.S)
    }
    
    sg := sig.Map[w]
    p := find(sg.S, stream)
    if p < 0 { return nil, nil, false }
    if p+len(sg.S)+8 > len(stream) { return nil, nil, false }
    
    wb := payload.BitsToBytes(stream[p+len(sg.S) : p+len(sg.S)+8])
    if len(wb) != 1 || int(wb[0]-'0') != w { return nil, nil, false }
    
    mb := payload.BitsToBytes(stream[p+len(sg.S)+8:])
    h, ok := meta.Unpack(mb)
    if !ok { return nil, nil, false }
    
    metaLen := 4 + 1 + 1 + 1 + 1 + len(h.Name) + 1 + len(h.Ext) + 8
    if metaLen > len(mb) { return nil, nil, false }
    
    pay := mb[metaLen:]
    if int(h.Size) > len(pay) { return nil, nil, false }
    pay = pay[:h.Size]
    
    return pay, &h, true
}

// DecodeFile decodes a steganographic MP3 file and extracts the hidden payload
func DecodeFile(inputFile, outputFile, key string, random, debug bool) error {
    // Read input file
    b, err := os.ReadFile(inputFile)
    if err != nil {
        return fmt.Errorf("failed to read input file: %v", err)
    }
    
    // Parse MP3
    f, err := mp3.Parse(b)
    if err != nil {
        return fmt.Errorf("failed to parse MP3: %v", err)
    }
    
    // Extract audio data
    var audio []byte
    for _, fr := range f.Frames {
        audio = append(audio, fr.Data...)
    }
    
    // Try different combinations of width and randomization
    for _, w := range []int{1, 2, 4} {
        for _, rnd := range []bool{random, !random} {
            pay, h, ok := tryDecode(audio, key, rnd, w, debug)
            if !ok {
                continue
            }
            
            // Decrypt if encrypted
            if (h.Flags & meta.FlagEncrypted) != 0 {
                pay = service.NewExtendedVigenereCipher(key).Decrypt(pay)
            }
            
            // Determine output filename
            fname := h.Name
            if outputFile != "" {
                fname = outputFile
            }
            
            // Write output file
            if err := os.WriteFile(fname, pay, 0644); err != nil {
                return fmt.Errorf("failed to write output file: %v", err)
            }
            
            fmt.Printf("Successfully decoded: width=%d bytes=%d file=%s\n", w, len(pay), fname)
            return nil
        }
    }
    
    return fmt.Errorf("signature not found - no hidden data detected")
}

func main() {
    // Example usage with manual variables
    inputFile := "stego.mp3"     // Input steganographic MP3 file
    outputFile := "decoded.txt"  // Output file (can be empty to use original name)
    key := "STEGANO"            // Encryption key/seed
    random := true              // Use random order
    debug := false              // Enable debug logs
    
    if err := DecodeFile(inputFile, outputFile, key, random, debug); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}