package processors

import (
	"testing"
	"processor"
	"storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
)



func TestDeliverProcessor(t *testing.T){




}

type testprocessor struct {
	*processor.GenericProcessor
	c chan types.App
	t *testing.T
}

func (tp testprocessor) Process(){
		app := types.GetApp("testid", tp.GetAgent())

		tp.c <- app


}
