package dock

import "github.com/joernweissenborn/aursirrt/core"


//A Docker handles app all communication.
type Docker interface {
	Launch(chan core.AppMessage, chan registerDockedApp) //Launches a docker
}

