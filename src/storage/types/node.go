package types

import (
	"github.com/joernweissenborn/aursir4go/messages"
	"dock/connection"
)


type Node interface {
	Exists() bool
	GetConnection() connection.Connection
	Create(DockMessage messages.DockMessage, Connection connection.Connection) bool
	GetExports() (exports []Export)
	GetImports() (imports []Import)
    Remove()
	IsApp()
}

