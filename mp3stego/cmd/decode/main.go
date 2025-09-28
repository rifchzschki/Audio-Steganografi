package main

import (
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"math/rand"
	"os"

	"mp3stego/internal/crypto"
	"mp3stego/internal/meta"
	"mp3stego/internal/mp3"
	"mp3stego/internal/payload"
	"mp3stego/internal/sig"
)

func seedFromKey(key string) int64 { h:=sha256.Sum256([]byte(key)); return int64(binary.LittleEndian.Uint64(h[:8])) }

func find(sigBits, hay []uint8) int { if len(sigBits)==0||len(hay)<len(sigBits){return -1}; for i:=0;i+len(sigBits)<=len(hay);i++{ ok:=true; for j:=0;j<len(sigBits);j++{ if hay[i+j]!=sigBits[j]{ ok=false; break } }; if ok{ return i } }; return -1 }

var dbg bool

func tryDecode(audio []byte, key string, random bool, w int) ([]byte, *meta.Header, bool) {
	eligible := make([]int,0,len(audio))
	for i,b:=range audio { if b!=0x00 && b!=0xFF { eligible=append(eligible,i) } }
	if len(eligible)==0 { return nil,nil,false }
	order := make([]int, len(audio))
	for i := range audio { order[i] = i }
	if len(order)==0 { panic("no audio bytes") }
	if random { rsrc:=rand.New(rand.NewSource(seedFromKey(key))); for i:=len(order)-1;i>0;i--{ j:=rsrc.Intn(i+1); order[i],order[j]=order[j],order[i] } }
	mask:=byte((1<<uint(w))-1)
	stream:=make([]uint8,0,len(order)*w)
	for _,pos:= range order { v:=audio[pos] & mask; for i:=w-1; i>=0; i-- { stream = append(stream, (v>>uint(i))&1) } }
	if dbg {
		nshow := 48
		if len(stream) < nshow { nshow = len(stream) }
		fmt.Printf("[DBG] w=%d random=%v firstBits=%v...\n", w, random, stream[:nshow])
		npos := 10
		if len(order) < npos { npos = len(order) }
		fmt.Printf("[DBG] firstPos=%v\n", order[:npos])
		sg:=sig.Map[w]
		fmt.Printf("[DBG] sigS=%v\n", sg.S)
	}
	sg:=sig.Map[w]
	p:=find(sg.S,stream); if p<0 { return nil,nil,false }
	if p+len(sg.S)+8>len(stream){ return nil,nil,false }
	wb:=payload.BitsToBytes(stream[p+len(sg.S):p+len(sg.S)+8]); if len(wb)!=1 || int(wb[0]-'0')!=w { return nil,nil,false }
	mb:=payload.BitsToBytes(stream[p+len(sg.S)+8:])
	h,ok:=meta.Unpack(mb); if !ok { return nil,nil,false }
	metaLen:=4+1+1+1+1+len(h.Name)+1+len(h.Ext)+8
	if metaLen>len(mb) { return nil,nil,false }
	pay:=mb[metaLen:]
	if int(h.Size) > len(pay) { return nil,nil,false }
	pay = pay[:h.Size]
	return pay,&h,true
}

func main(){
	in:=flag.String("in","","stego mp3")
	out:=flag.String("out","","output file")
	key:=flag.String("key","STEGANO","key/seed")
	random:=flag.Bool("random",true,"random order")
	flag.BoolVar(&dbg, "dbg", false, "debug logs")
	flag.Parse()
	if *in=="" { flag.Usage(); os.Exit(2) }
	b,err:=os.ReadFile(*in); if err!=nil{ panic(err) }
	f,err:=mp3.Parse(b); if err!=nil{ panic(err) }
	var audio []byte
	for _,fr:=range f.Frames{ audio=append(audio, fr.Data...) }
	for _,w:=range []int{1,2,4} {
		for _,rnd := range []bool{*random, !*random} {
			pay,h,ok := tryDecode(audio, *key, rnd, w)
			if !ok { continue }
			if (h.Flags & meta.FlagEncrypted) != 0 { pay = crypto.NewExtendedVigenere(*key).Decrypt(pay) }
			fname:=h.Name
			if *out!="" { fname = *out }
			if err:=os.WriteFile(fname,pay,0644); err!=nil{ panic(err) }
			fmt.Printf("ok width=%d bytes=%d file=%s\n", w, len(pay), fname)
			return
		}
	}
	panic("signature not found")
}
