package dockprocessor

import (
	"testing"
	"processor"
	"storage"
	"github.com/joernweissenborn/aursir4go"
	"storage/types"
)

var Testdockmsg = aursir4go.AurSirDockMessage{"testapp",[]string{"JSON"}}


func TestDockProcessor(t *testing.T){

	pc := processor.Testprocessor()
	defer close(pc)

	var dp DockProcessor
	dp.GenericProcessor = processor.GetGenericProcessor()
	dp.AppId = "testid"
	dp.DockMessage = Testdockmsg

	pc <- dp

	var tp testprocessor
	tp.GenericProcessor = processor.GetGenericProcessor()
	tp.t = t
	tp.c = make(chan types.App)
	defer close(tp.c)

	pc <- tp
	exists := (<-tp.c).Exists()
	if !exists{
		tp.t.Error("App does not exist")
	}

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
