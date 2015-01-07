package addimportprocessor

import (
	"testing"
	"processor"
	"github.com/joernweissenborn/aursir4go"
	"storage/types"
	"processor/processors/dockprocessor"
)

var Testdockmsg = aursir4go.AurSirDockMessage{"testapp",[]string{"JSON"}}
var Testaddimpmsg = aursir4go.AurSirAddImportMessage{aursir4go.HelloAurSirAppKey, []string{"one","two"}}


func TestAddExportProcessor(t *testing.T){

	pc := processor.Testprocessor()
	defer close(pc)

	var dp dockprocessor.DockProcessor
	dp.GenericProcessor = processor.GetGenericProcessor()
	dp.AppId = "testid"
	dp.DockMessage = Testdockmsg

	go func (){pc <- dp}()

	var ap AddImportProcessor
	ap.GenericProcessor = processor.GetGenericProcessor()
	ap.AppId = "testid"
	ap.AddImportMsg = Testaddimpmsg

	go func (){pc <- ap}()

	var tp testprocessor
	tp.GenericProcessor = processor.GetGenericProcessor()
	tp.c = make(chan bool)
	defer close(tp.c)

	go func (){pc <- tp}()
	exp := <- tp.c
	if exp {
		t.Error("Failed to create import")
	}


}

type testprocessor struct {
	*processor.GenericProcessor
	c chan bool
	t *testing.T
}

func (tp testprocessor) Process(){
		Import := types.GetImport("testid",aursir4go.HelloAurSirAppKey,[]string{"one","two"},tp.GetAgent())
		tp.c <- Import.GetId()==""
}