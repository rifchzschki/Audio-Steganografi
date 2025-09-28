package sig

type Pair struct{ 
	S, E []uint8 
}

func Bits(s string) []uint8 { 
	r := make([]uint8, len(s))
	for i, c := range s { 
		if c=='1' { 
			r[i]=1 
		} 
	} 
	return r 
}

var Map = map[int]Pair{ 
	1:{Bits("10101010101010"),Bits("10101010101010")}, 
	2:{Bits("01010101010101"),Bits("01010101010101")}, 
	3:{Bits("10101010101010"),Bits("01010101010101")},
	4:{Bits("01010101010101"),Bits("10101010101010")}, 
}

func WidthByte(n int) []uint8 { 
	x:=byte('0'+byte(n))
	r:=make([]uint8,8)
	for i:=7;i>=0;i--{ 
		if (x>>uint(i))&1==1 { 
			r[7-i]=1 
		} 
	} 
	return r 
}
