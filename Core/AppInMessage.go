package core
import (
	"github.com/joernweissenborn/aursir4go"
)


//An AppInMessage is used to send incoming messages to the core
type AppMessage struct {

	SenderUUID string //ID of the sending app

	AppMsg aursir4go.AppMessage

}
