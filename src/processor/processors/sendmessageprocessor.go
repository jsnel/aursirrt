package processors

import (
	"processor"
	"storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
	"github.com/joernweissenborn/aursir4go/util"
	"dock/connection"
)

type SendMessageProcessor struct {

	*processor.GenericProcessor

	App types.App

	Msg messages.AurSirMessage
}

func (p SendMessageProcessor) Process() {
	               //printDebug(fmt.Sprint("sendmsg",p.Msg) )
	conn := p.App.GetConnection()
	id := p.App.Id

	var msgtype int64
	var codec string = "JSON"

	switch p.Msg.(type) {

	case messages.DockedMessage:
		msgtype = messages.DOCKED

	case messages.ExportAddedMessage:
		msgtype = messages.EXPORT_ADDED


	case messages.ImportAddedMessage:
		msgtype = messages.IMPORT_ADDED

	case messages.ImportUpdatedMessage:
		msgtype = messages.IMPORT_UPDATED
	case *messages.Request:
		msgtype = messages.REQUEST

case messages.Result:
		msgtype = messages.RESULT

	}

	decoder := util.GetCodec(codec)
	if decoder == nil {
		return
	}
	msg, err := decoder.Encode(p.Msg)
	if err != nil {
		return
	}
	go p.send(id,conn,msgtype,codec,msg)

}

func (p SendMessageProcessor) send(id string,conn connection.Connection,msgtype int64, codec string,msg []byte) {

	if nil != conn.Send(msgtype,codec,msg) {
		var lp LeaveProcessor
		lp.AppId = id
		p.SpawnProcess(lp)
	}
}

