package aursirscript

import "unicode/utf8"


func IsComment(code []byte) (IsComment bool, err error ){

	rune, size := utf8.DecodeRune(code)
	if rune == '/'{
		nextrune,_ := utf8.DecodeRune(code[size:])
		if nextrune == '/' {
			IsComment = true
		} else {
			err = CommentError{}
		}
	}

	return
}
