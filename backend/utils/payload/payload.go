package payload

func ToBits(b []byte) []uint8 { 
	out:=make([]uint8,0,len(b)*8)
	for _,v:=range b{ 
		for i:=7;i>=0;i--{ 
			if (v>>uint(i))&1==1{ 
				out=append(out,1)
			}else{ 
				out=append(out,0)
			}
		}
	} 
	return out 
}

func BitsToBytes(bits []uint8) []byte { 
	n:=len(bits)/8
	out:=make([]byte,0,n)
	for i:=0;i+8<=len(bits);i+=8{ 
		var v byte
		for j:=0;j<8;j++{ 
			v=(v<<1)|bits[i+j] 
		} 
		out=append(out,v) 
	} 
	return out 
}
