package mp3

import (
	"bytes"
	"errors"
)

type FrameHeader struct {
    VersionID   int
    Layer       int
    Protection  bool
    Bitrate     int
    SampleRate  int
    Padding     bool
    ChannelMode int
    FrameLength int
}

type Frame struct {
    Header      *FrameHeader
    HeaderBytes []byte
    Data        []byte
}

type File struct {
    ID3v2  []byte
    Frames []*Frame
    ID3v1  []byte
}

func Parse(data []byte) (*File, error) {
    f := &File{}
    i := 0
    
    if len(data) >= 10 && bytes.Equal(data[0:3], []byte("ID3")) {
        size := 0
        for j := 0; j < 4; j++ {
            size = (size << 7) | int(data[6+j]&0x7F)
        }
        if 10+size <= len(data) {
            f.ID3v2 = make([]byte, 10+size)
            copy(f.ID3v2, data[:10+size])
            i = 10 + size
        }
    }
    
    if len(data) >= 128 && bytes.Equal(data[len(data)-128:len(data)-125], []byte("TAG")) {
        f.ID3v1 = make([]byte, 128)
        copy(f.ID3v1, data[len(data)-128:])
    }
    
    max := len(data)
    if len(f.ID3v1) > 0 {
        max -= 128
    }
    
    // Parse MP3 frames
    for i+4 <= max {
        if data[i] != 0xFF || (data[i+1]&0xE0) != 0xE0 {
            i++
            continue
        }
        
        // Calculate frame length
        fl := frameLen(data[i:])
        if fl <= 0 || i+fl > max {
            i++
            continue
        }
        
        // Create frame header bytes
        fh := make([]byte, 4)
        copy(fh, data[i:i+4])
        
        // Create frame with data
        fr := &Frame{
            Header:      &FrameHeader{FrameLength: fl},
            HeaderBytes: fh,
            Data:        append([]byte(nil), data[i+4:i+fl]...),
        }
        
        f.Frames = append(f.Frames, fr)
        i += fl
    }
    
    if len(f.Frames) == 0 {
        return nil, errors.New("no frames")
    }
    
    return f, nil
}


func frameLen(h []byte) int {
    if len(h) < 4 {
        return -1
    }
    
    brIdx := int((h[2] >> 4) & 0x0F)
    srIdx := int((h[2] >> 2) & 0x03)
    
    if brIdx == 0 || brIdx == 15 || srIdx == 3 {
        return -1
    }
    
    pad := (h[2] >> 1) & 1
    ver := (h[1] >> 3) & 3
    lay := (h[1] >> 1) & 3
    
    if lay != 1 {
        return -1
    }
    
    // Bitrate table (Layer III)
    brTab := []int{0, 32000, 40000, 48000, 56000, 64000, 80000, 96000, 
        112000, 128000, 160000, 192000, 224000, 256000, 320000}
    
    // Sample rate table
    srTab := []int{44100, 48000, 32000}
    
    sr := srTab[srIdx]
    if ver == 2 {
        sr /= 2
    } else if ver == 0 {
        sr /= 4
    }
    
    if sr == 0 {
        return -1
    }
    
    return (144 * brTab[brIdx] / sr) + int(pad)
}

func Serialize(f *File) []byte {
    var out []byte
    
    if len(f.ID3v2) > 0 {
        out = append(out, f.ID3v2...)
    }
    
    for _, fr := range f.Frames {
        out = append(out, fr.HeaderBytes...)
        out = append(out, fr.Data...)
    }
    
    if len(f.ID3v1) > 0 {
        out = append(out, f.ID3v1...)
    }
    
    return out
}
