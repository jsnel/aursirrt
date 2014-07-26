package Dock

import "github.com/joernweissenborn/AurSirRt/Core"


//A Docker handles app all communication.
type Docker interface {
	Launch(chan Core.AppMessage, chan registerDockedApp) //Launches a docker
}

