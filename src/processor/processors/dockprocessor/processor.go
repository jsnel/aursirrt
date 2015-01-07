package dockprocessor

import (
	"processor"
	"storage/types"
	"dock"
	"github.com/joernweissenborn/aursir4go/messages"
)

type DockProcessor struct {

	*processor.GenericProcessor

	AppId string

	DockMessage messages.DockMessage

	Connection dock.Connection

}

func (p DockProcessor) Process() {
	if p.Connection != nil {

		app := types.GetApp(p.AppId, p.GetAgent())
		ok := app.Create(p.DockMessage, p.Connection)
		app.GetConnection().Send(messages.DockedMessage{ok})

	}

}

