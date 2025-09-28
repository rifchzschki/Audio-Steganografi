package utils

// import(
// 	"os"
// 	"github.com/hajimehoshi/go-mp3"
// )

// type Sig struct{ Start, End []uint8 }

// func bits(s string) []uint8 {
// 	b := make([]uint8, len(s))
// 	for i, c := range s {
// 		if c == '1' {
// 			b[i] = 1
// 		} else {
// 			b[i] = 0
// 		}
// 	}
// 	return b
// }

// var signature = map[int]Sig{
// 	1: {Start: bits("10101010101010"), End: bits("10101010101010")},
// 	2: {Start: bits("01010101010101"), End: bits("01010101010101")},
// 	3: {Start: bits("10101010101010"), End: bits("01010101010101")},
// 	4: {Start: bits("01010101010101"), End: bits("10101010101010")},
// }