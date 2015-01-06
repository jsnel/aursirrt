package types

import (
	"testing"
	"storage"
	"github.com/joernweissenborn/aursir4go"
)

func TestExportCreation(t *testing.T){
	agent := storage.NewAgent()
	app := GetApp("testid",agent)
	dockmsg := aursir4go.AurSirDockMessage{"HelloWorld",[]string{"JSON"}}
	app.Create(dockmsg)

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
}
