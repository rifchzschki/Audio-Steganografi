package crypto

type ExtendedVigenere struct{ Key []byte }

func NewExtendedVigenere(key string) *ExtendedVigenere { return &ExtendedVigenere{Key: []byte(key)} }

func (e *ExtendedVigenere) Encrypt(p []byte) []byte {
	if len(e.Key) == 0 { out := make([]byte, len(p)); copy(out, p); return out }
	out := make([]byte, len(p))
	for i, b := range p { out[i] = byte(int(b)+int(e.Key[i%len(e.Key)])&0xFF) }
	return out
}

func (e *ExtendedVigenere) Decrypt(c []byte) []byte {
	if len(e.Key) == 0 { out := make([]byte, len(c)); copy(out, c); return out }
	out := make([]byte, len(c))
	for i, b := range c { out[i] = byte(int(b)-int(e.Key[i%len(e.Key)])&0xFF) }
	return out
}
