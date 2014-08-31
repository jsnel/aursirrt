package aursirscript

import "unicode/utf8"

func GetLiteral(code []byte) (Literal []byte, RemainingCode []byte){
	RemainingCode =code
	cont := true
	Literal = []byte{}
	for len(RemainingCode) != 0 && cont{
		rune, size := utf8.DecodeRune(code)

		if rune != ' ' {
			for _,char := range RemainingCode[:size] {
				Literal = append(Literal,char)
			}
		} else {
			cont = false
		}
		RemainingCode = RemainingCode[size:]

	}

	return
}
