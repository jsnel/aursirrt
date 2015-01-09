package types
import (
	"testing"
	"storage"
	"github.com/joernweissenborn/aursir4go"
	"github.com/joernweissenborn/aursir4go/messages"
)

func TestImport(t *testing.T){
	agent := storage.NewAgent()
	app := GetApp("testid",agent)
	dockmsg := messages.DockMessage{"HelloWorld",[]string{"JSON"}}
	app.Create(dockmsg,testconn{})

	Import := GetImport("",aursir4go.HelloAurSirAppKey, []string{"one","two"},agent)
	Import.Add()

	if Import.GetId() != "" {
		t.Error("Created Import for non existing app")
	}
	Import = GetImport("testid",aursir4go.HelloAurSirAppKey, []string{"one","two"},agent)
	Import.Add()

	if Import.GetId() == "" {
		t.Error("Could not add Import")
	}
	Import = GetImport("testid",aursir4go.HelloAurSirAppKey, []string{"one","two"},agent)

	if Import.GetId() == "" {
		t.Error("Could not retrieve Import")
	}
	key := Import.GetAppKey()
	if !key.Exists(){
		t.Error("Could not create key")
	}

	if len(key.GetImporter()) == 0 {
		t.Error("Could not retrieve import from key")

	}

	if Import.HasExporter() {
		t.Error("Exporter should not be present")
	}
	eapp := GetApp("testexp",agent)
	eapp.Create(dockmsg,testconn{})

	export := GetExport("testexp",aursir4go.HelloAurSirAppKey, []string{"one","two"},agent)
	export.Add()
	if !Import.HasExporter() {
		t.Error("Exporter should be present")
	}
}
