package dock

import (
	"github.com/joernweissenborn/aursir4go/messages"
)


//A Docker handles app all communication.
type Docker interface {
	//Launch(chan core.AppMessage, chan interface {}) //Launches a docker
}

type Connection interface {
	Send(msg messages.AurSirMessage)
}
