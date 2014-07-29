package core
import (
	"github.com/joernweissenborn/AurSir4Go"
)


//An AppInMessage is used to send incoming messages to the core
type AppMessage struct {

	SenderUUID string //ID of the sending app

	AppMsg AurSir4Go.AppMessage

}
