package types

import (
	"testing"
	"storage"
	"github.com/joernweissenborn/aursir4go"
)

func TestAppCreation(t *testing.T){
	agent := storage.NewAgent()
	app := GetApp("testid",agent)
	if app.Exists() {
		t.Error("Found non existing app")
	}
	dockmsg := aursir4go.AurSirDockMessage{"HelloWorld",[]string{"JSON"}}
	app.Create(dockmsg)
	if !app.Exists() {
		t.Error("Could not create app")
	}
}

func TestAppRemoval(t *testing.T){
	agent := storage.NewAgent()
	app := GetApp("testid",agent)
	dockmsg := aursir4go.AurSirDockMessage{"HelloWorld",[]string{"JSON"}}
	app.Create(dockmsg)
	app.Remove()
	if app.Exists() {
		t.Error("Could not remove app")
	}
}
