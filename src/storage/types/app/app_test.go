package app

import (
	"testing"
	"storage"
	"github.com/joernweissenborn/aursir4go"
)

func TestAppCreation(t *testing.T){
	agent := storage.NewAgent()
	app := Get("testid",agent)
	dockmsg := aursir4go.AurSirDockMessage{"HelloWorld",[]string{"JSON"}}
	app.Create(dockmsg)
	if !app.Exists() {
		t.Error("Could not create app")
	}
}
