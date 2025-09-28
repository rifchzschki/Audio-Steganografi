package meta

import "encoding/binary"

const Magic = 0x6D703373

type Flags uint8

const ( FlagEncrypted Flags = 1<<0; FlagRandomStart Flags = 1<<1 )

type Header struct{ Version uint8; Flags Flags; NLSB uint8; Name string; Size uint64; Ext string }

func Pack(h Header) []byte { name:=[]byte(h.Name); ext:=[]byte(h.Ext); b:=make([]byte,0,4+1+1+1+1+len(name)+1+len(ext)+8); tmp:=make([]byte,8); binary.BigEndian.PutUint32(tmp[:4],Magic); b=append(b,tmp[:4]...); b=append(b,h.Version); b=append(b,byte(h.Flags)); b=append(b,byte(h.NLSB)); b=append(b,byte(len(name))); b=append(b,name...); b=append(b,byte(len(ext))); b=append(b,ext...); binary.BigEndian.PutUint64(tmp,h.Size); b=append(b,tmp...); return b }

func Unpack(b []byte) (Header,bool) { var h Header; if len(b)<4+1+1+1+1+1+8 { return h,false }; if binary.BigEndian.Uint32(b[:4])!=Magic { return h,false }; i:=4; h.Version=b[i]; i++; h.Flags=Flags(b[i]); i++; h.NLSB=b[i]; i++; nl:=int(b[i]); i++; if i+nl>len(b){ return h,false }; h.Name=string(b[i:i+nl]); i+=nl; el:=int(b[i]); i++; if i+el>len(b){ return h,false }; h.Ext=string(b[i:i+el]); i+=el; if i+8>len(b){ return h,false }; h.Size=binary.BigEndian.Uint64(b[i:i+8]); return h,true }
