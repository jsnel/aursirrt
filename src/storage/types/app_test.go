package types

import (
	"testing"
	"storage"
	"github.com/joernweissenborn/aursir4go/messages"
)


func TestAppCreation(t *testing.T){
	agent := storage.NewAgent()
	app := GetApp("testid",agent)
	if app.Exists() {
		t.Error("Found non existing app")
	}
	dockmsg := messages.DockMessage{"HelloWorld",[]string{"JSON"}}
	if !app.Create(dockmsg,testconn{}) {
		t.Error("Could not create app")
	}
	if !app.Exists() {
		t.Error("Could not find app")
	}
	if app.Create(dockmsg,testconn{}) {
		t.Error("Could create app")
	}
}

func TestAppRemoval(t *testing.T){
	agent := storage.NewAgent()
	app := GetApp("testid",agent)
	dockmsg := messages.DockMessage{"HelloWorld",[]string{"JSON"}}
	app.Create(dockmsg,testconn{})
	app.Remove()
	if app.Exists() {
		t.Error("Could not remove app")
	}
}

type testconn struct {

}

func (testconn) Send(msgtype int64, codec string,msg []byte) (err error) {
	return
}
func (testconn) Init()error{
	return nil
}
