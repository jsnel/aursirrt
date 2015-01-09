package types

import (
	"testing"
	"storage"
	"github.com/joernweissenborn/aursir4go"
	"github.com/joernweissenborn/aursir4go/messages"
)

func TestExport(t *testing.T){
	agent := storage.NewAgent()
	app := GetApp("testid",agent)
	dockmsg := messages.DockMessage{"HelloWorld",[]string{"JSON"}}
	app.Create(dockmsg,testconn{})

	export := GetExport("",aursir4go.HelloAurSirAppKey, []string{"one","two"},agent)
	export.Add()

	if export.GetId() != "" {
		t.Error("Created export for non existing app")
	}
	export = GetExport("testid",aursir4go.HelloAurSirAppKey, []string{"one","two"},agent)
	export.Add()

	if export.GetId() == "" {
		t.Error("Could not add export")
	}
	export = GetExport("testid",aursir4go.HelloAurSirAppKey, []string{"one","two"},agent)
	export.Add()

	if export.GetId() == "" {
		t.Error("Could not retrieve export")
	}

	if len(export.GetAppKey().GetExporter()) == 0 {
		t.Error("Could not retrieve export from key")

	}

}
