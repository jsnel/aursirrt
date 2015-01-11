package types

import (
	"testing"
	"storage"
	"github.com/joernweissenborn/aursir4go/messages"
	"github.com/joernweissenborn/aursir4go/Example/keys"
)

func TestExport(t *testing.T){
	agent := storage.NewAgent()
	app := GetApp("testid",agent)
	dockmsg := messages.DockMessage{"HelloWorld",[]string{"JSON"},false}
	app.Create(dockmsg,testconn{})

	export := GetExport("",keys.HelloAurSirAppKey, []string{"one","two"},"",agent)
	export.Add()

	if export.GetId() != "" {
		t.Error("Created export for non existing app")
	}
	export = GetExport("testid",keys.HelloAurSirAppKey, []string{"one","two"},"",agent)
	export.Add()

	if export.GetId() == "" {
		t.Error("Could not add export")
	}
	export = GetExport("testid",keys.HelloAurSirAppKey, []string{"one","two"},"",agent)
	export.Add()

	if export.GetId() == "" {
		t.Error("Could not retrieve export")
	}
	             appkey := export.GetAppKey()
	if len(appkey.GetExporter()) == 0 {
		t.Error("Could not retrieve export from key")

	}

}
