package deliverjobprocessor

import (
	"testing"
	"processor"
	"storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
	"processor/processors/dockprocessor"
)

var Testdockmsg = messages.DockMessage{"testapp",[]string{"JSON"}}


func TestDockProcessor(t *testing.T){

	pc := processor.Testprocessor()
	defer close(pc)

	c := make(chan bool)
	defer close(c)
	conn := testconnection{c}

	var dp dockprocessor.DockProcessor
	dp.GenericProcessor = processor.GetGenericProcessor()
	dp.AppId = "testid"
	dp.DockMessage = Testdockmsg
	dp.Connection = conn

	pc <- dp
	if !(<-c) {
		t.Error("Dockmessages says not ok")
	}



	var tp testprocessor
	tp.GenericProcessor = processor.GetGenericProcessor()
	tp.t = t
	tp.c = make(chan types.App)
	defer close(tp.c)

	pc <- tp
	exists := (<-tp.c).Exists()
	if !exists{
		t.Error("App does not exist")
	}

	pc <- dp
	if (<-c) {
		t.Error("Dockmessages says  ok")
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

type testconnection struct {
	c chan bool
}

func (tc testconnection) Send(msg messages.AurSirMessage) {

	tc.c <- msg.(messages.DockedMessage).Ok

}
