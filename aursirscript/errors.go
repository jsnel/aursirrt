package aursirscript

type CommentError struct {

}

func (CommentError) Error() string {
	return "ParsingError: Forbidden use of '/'"
}
