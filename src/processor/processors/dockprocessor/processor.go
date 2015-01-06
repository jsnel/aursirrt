package dockprocessor

import (
	"github.com/joernweissenborn/aursir4go"
	"processor"
	"storage/types"
)

type DockProcessor struct {

	*processor.GenericProcessor

	AppId string

	DockMessage aursir4go.AurSirDockMessage

}

func (p DockProcessor) Process() {

	app := types.GetApp(p.AppId,p.GetAgent())
	if !app.Exists() {
		app.Create(p.DockMessage)
	}

}

