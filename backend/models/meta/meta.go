package meta

import "encoding/binary"

const Magic = 0x6D703373

type Flags uint8

const (
    FlagEncrypted   Flags = 1 << 0 
    FlagRandomStart Flags = 1 << 1 
)

type Header struct {
    Version uint8  // Version of the steganographic format
    Flags   Flags  // Flags indicating encryption and other options
    NLSB    uint8  // Number of LSBs used
    Name    string // Original filename
    Size    uint64 // Size of the payload in bytes
    Ext     string // File extension
}

func Pack(h Header) []byte {
    name := []byte(h.Name)
    ext := []byte(h.Ext)
    
    // Calculate required buffer size
    bufSize := 4 + 1 + 1 + 1 + 1 + len(name) + 1 + len(ext) + 8
    b := make([]byte, 0, bufSize)
    tmp := make([]byte, 8)
    
    // Pack magic number
    binary.BigEndian.PutUint32(tmp[:4], Magic)
    b = append(b, tmp[:4]...)
    
    // Pack header fields
    b = append(b, h.Version)
    b = append(b, byte(h.Flags))
    b = append(b, byte(h.NLSB))
    
    // Pack name
    b = append(b, byte(len(name)))
    b = append(b, name...)
    
    // Pack extension
    b = append(b, byte(len(ext)))
    b = append(b, ext...)
    
    // Pack size
    binary.BigEndian.PutUint64(tmp, h.Size)
    b = append(b, tmp...)
    
    return b
}

func Unpack(b []byte) (Header, bool) {
    var h Header
    
    // Check minimum required length
    if len(b) < 4+1+1+1+1+1+8 {
        return h, false
    }
    
    // Check magic number
    if binary.BigEndian.Uint32(b[:4]) != Magic {
        return h, false
    }
    
    i := 4
    
    // Unpack header fields
    h.Version = b[i]
    i++
    h.Flags = Flags(b[i])
    i++
    h.NLSB = b[i]
    i++
    
    // Unpack name
    nl := int(b[i])
    i++
    if i+nl > len(b) {
        return h, false
    }
    h.Name = string(b[i : i+nl])
    i += nl
    
    // Unpack extension
    el := int(b[i])
    i++
    if i+el > len(b) {
        return h, false
    }
    h.Ext = string(b[i : i+el])
    i += el
    
    // Unpack size
    if i+8 > len(b) {
        return h, false
    }
    h.Size = binary.BigEndian.Uint64(b[i : i+8])
    
    return h, true
}
