package processors

import (
	"processor"
	"github.com/joernweissenborn/aursir4go/util"
	"github.com/joernweissenborn/aursir4go/messages"
)


type ParseMessageProccesor struct {

	*processor.GenericProcessor

	AppId string
	Codec string
	Type int64
	Msg []byte
}

func (p ParseMessageProccesor) Process() {

	decoder := util.GetCodec(p.Codec)
	if decoder == nil {
		return
	}

	switch p.Type {
	case messages.LEAVE:
		var lp LeaveProcessor
			lp.AppId == p.AppId
		p.SpawnProcess(lp)

	case messages.ADD_EXPORT:
		var m messages.AddExportMessage
		decoder.Decode(p.Msg, &m)
		var np AddExportProcessor
		np.AppId = p.AppId
		np.AddExportMsg = m
		p.SpawnProcess(np)

	case messages.UPDATE_EXPORT:
		var m messages.UpdateExportMessage
		decoder.Decode(p.Msg, &m)
		var np UpdateExportProcessor
		np.AppId = p.AppId
		np.UpdateExportMsg = m
		np.SpawnProcess(np)

	case messages.ADD_IMPORT:
		var m messages.AddImportMessage
		decoder.Decode(p.Msg, &m)
		var np AddImportProcessor
		np.AppId = p.AppId
		np.AddImportMsg = m
		np.SpawnProcess(np)

	case messages.UPDATE_IMPORT:
		var m messages.UpdateImportMessage
		decoder.Decode(p.Msg, &m)
		var np UpdateImportProcessor
		np.AppId = p.AppId
		np.UpdateImportMsg = m
		np.SpawnProcess(np)

	case messages.REQUEST:
		var m messages.Request
		decoder.Decode(p.Msg, &m)
		var np RequestProcessor
		np.AppId = p.AppId
		np.Request = m
		np.SpawnProcess(np)
	}

}

