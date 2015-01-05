package dockprocessor

import (
	"github.com/joernweissenborn/aursir4go"
	"processor"
)

type DockProcessor struct {

	processor.GenericProcessor

	AppId string

	DockMessage aursir4go.AurSirDockMessage

}

func (p DockProcessor) Process() {


	if !app.Exists("") {
		app.Create(p.AppId,p.DockMessage)
	}

}

func registerApp(){}
