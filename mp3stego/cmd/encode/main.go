package main

import (
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"mp3stego/internal/crypto"
	"mp3stego/internal/meta"
	"mp3stego/internal/mp3"
	"mp3stego/internal/payload"
	"mp3stego/internal/sig"
)

func seedFromKey(key string) int64 { h:=sha256.Sum256([]byte(key)); return int64(binary.LittleEndian.Uint64(h[:8])) }

func main(){
	in:=flag.String("in","","input mp3")
	secret:=flag.String("secret","","secret file")
	out:=flag.String("out","stego.mp3","output")
	n:=flag.Int("n",1,"1|2|4")
	key:=flag.String("key","STEGANO","key/seed")
	enc:=flag.Bool("enc",false,"encrypt")
	random:=flag.Bool("random",true,"random order")
	flag.Parse()
	if *in==""||*secret=="" { flag.Usage(); os.Exit(2) }
	if *n!=1&&*n!=2&&*n!=4 { panic("n must 1,2,4") }
	coverBytes,err:=os.ReadFile(*in); if err!=nil{ panic(err) }
	f,err:=mp3.Parse(coverBytes); if err!=nil{ panic(err) }
	name:=filepath.Base(*secret); ext:=filepath.Ext(name)
	sb,err:=os.ReadFile(*secret); if err!=nil{ panic(err) }
	if *enc { sb = crypto.NewExtendedVigenere(*key).Encrypt(sb) }
	h:=meta.Header{Version:1, Flags:0, NLSB:uint8(*n), Name:name, Size:uint64(len(sb)), Ext:ext}
	if *enc { h.Flags|=meta.FlagEncrypted }
	if *random { h.Flags|=meta.FlagRandomStart }
	metaBytes:=meta.Pack(h)
	S:=sig.Map[*n]
	bits:=make([]uint8,0)
	bits=append(bits,S.S...)
	bits=append(bits,sig.WidthByte(*n)...)
	bits=append(bits,payload.ToBits(metaBytes)...)
	bits=append(bits,payload.ToBits(sb)...)
	bits=append(bits,S.E...)
	var audio []byte
	for _,fr:=range f.Frames{ audio=append(audio, fr.Data...) }
	order := make([]int, len(audio))
	for i := range audio { order[i] = i }
	if len(order)==0 { panic("no audio bytes") }
	capBits := len(order) * (*n)
	if capBits < len(bits) { panic("capacity too small") }
	if *random { rsrc:=rand.New(rand.NewSource(seedFromKey(*key))); for i:=len(order)-1; i>0; i-- { j:=rsrc.Intn(i+1); order[i],order[j]=order[j],order[i] } }
	mask := byte((1<<uint(*n)) - 1)
	bi := 0
	need := (len(bits)+*n-1)/(*n)
	for t:=0; t<need; t++ {
		pos := order[t]
		var pv byte = 0
		for i:=0; i<*n && bi<len(bits); i++ { pv = (pv<<1) | bits[bi]; bi++ }
		if bi%(*n) != 0 { pv <<= uint((*n) - (bi%(*n))) }
		audio[pos] = (audio[pos] &^ mask) | (pv & mask)
	}
	idx:=0
	for _,fr:=range f.Frames { for k:=range fr.Data { if idx<len(audio){ fr.Data[k]=audio[idx]; idx++ } } }
	outBytes:=mp3.Serialize(f)
	if err:=os.WriteFile(*out,outBytes,0644); err!=nil{ panic(err) }
	fmt.Printf("ok bits=%d n=%d\n", len(bits), *n)
}
