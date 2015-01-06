package addexportprocessor

import (
	"testing"
	"processor"
	"github.com/joernweissenborn/aursir4go"
	"storage/types"
	"processor/processors/dockprocessor"
)

var Testdockmsg = aursir4go.AurSirDockMessage{"testapp",[]string{"JSON"}}
var Testaddexpmsg = aursir4go.AurSirAddExportMessage{aursir4go.HelloAurSirAppKey, []string{"one","two"}}


func TestAddExportProcessor(t *testing.T){

	pc := processor.Testprocessor()
	defer close(pc)

	var dp dockprocessor.DockProcessor
	dp.GenericProcessor = processor.GetGenericProcessor()
	dp.AppId = "testid"
	dp.DockMessage = Testdockmsg

	pc <- dp

	var ap AddExportProcessor
	ap.GenericProcessor = processor.GetGenericProcessor()
	ap.AppId = "testid"
	ap.AddExportMsg = Testaddexpmsg

	pc <- ap

	var tp testprocessor
	tp.GenericProcessor = processor.GetGenericProcessor()
	tp.c = make(chan *types.Export)
	defer close(tp.c)

	pc <- tp
	 <- tp.c
//	if exp.GetId() == "" {
//		t.Error("Failed to create export")
//	}


}

type testprocessor struct {
	*processor.GenericProcessor
	c chan *types.Export
	t *testing.T
}

func (tp testprocessor) Process(){
		//export := types.GetExport("testid",aursir4go.HelloAurSirAppKey,[]string{"one","two"},tp.GetAgent())
		tp.c <- new(types.Export)
}
