/*
	Package Core is responsible for processing incoming messages
*/
package Core

import (
	"github.com/joernweissenborn/AurSir4Go"
	"github.com/joernweissenborn/AurSirRt/StorageCore"
	"log"
)

func Launch(AppInChannel, AppOutChannel chan AppMessage) {

	log.Println("Core Launching")

	var c core

	c.storageAgent.Launch()

	c.appOutChannel = AppOutChannel

	go c.routeIncomingAppMsg(AppInChannel)

}

type core struct {
	storageAgent StorageCore.StorageCoreAgent

	appOutChannel chan AppMessage
}

func (c core) routeIncomingAppMsg(appInChannel chan AppMessage) {

	for AppMessage := range appInChannel {
		aursirMessage, err := AppMessage.AppMsg.Decode()

		if err == nil {

			switch aursirMessage := aursirMessage.(type) {

			case AurSir4Go.AurSirDockMessage:
				go c.dock(AppMessage.SenderUUID, aursirMessage)

			case AurSir4Go.AurSirLeaveMessage:
				go c.leave(AppMessage.SenderUUID)

			case AurSir4Go.AurSirAddExportMessage:
				go c.addExport(AppMessage.SenderUUID, aursirMessage)

			case AurSir4Go.AurSirAddImportMessage:
				go c.addImport(AppMessage.SenderUUID, aursirMessage)

			case AurSir4Go.AurSirRequest:
				go c.request(AppMessage.SenderUUID,aursirMessage)

			case AurSir4Go.AurSirResult:
				go c.result(AppMessage.SenderUUID,aursirMessage)
			case AurSir4Go.AurSirListenMessage:
				go c.listen(AppMessage.SenderUUID,aursirMessage)
			case AurSir4Go.AurSirStopListenMessage:
				go c.stopListen(AppMessage.SenderUUID,aursirMessage)

			default:
				log.Println("unknown message")
			}
		}
	}
}

func (c core) listen(senderId string, lmsg AurSir4Go.AurSirListenMessage){
	log.Println("Processing LISTEN request from", senderId)
	c.storageAgent.Write(StorageCore.ListenRequest{senderId,lmsg.FunctionName,lmsg.ImportId})
}

func (c core) stopListen(senderId string, slmsg AurSir4Go.AurSirStopListenMessage){
	log.Println("Processing STOP_LISTEN request from", senderId)

}

func (c core) result(senderId string, rmsg AurSir4Go.AurSirResult) {
	log.Println("Processing RESULT request from", senderId)
	reply := c.storageAgent.Write(StorageCore.AddResRequest{senderId,rmsg})
	reqreg, ok := reply.(StorageCore.ResRegistered)
	log.Println("StorageCore error ",reqreg)
	if ok && len(reqreg.Importer)!=0{
		for _, imp := range reqreg.Importer{
			var rm AurSir4Go.AppMessage
			rm.Encode(rmsg,"JSON")
			c.appOutChannel <-AppMessage{imp,rm}
		}
	}
}

func (c core) request(senderId string, rmsg AurSir4Go.AurSirRequest) {
	log.Println("Processing REQUEST request from", senderId)
	reply := c.storageAgent.Write(StorageCore.AddReqRequest{senderId,rmsg})
	log.Println(reply)
	reqreg, ok := reply.(StorageCore.ReqRegistered)
	log.Println(reqreg)
	if ok && len(reqreg.Exporter)!=0{
		for _, exp := range reqreg.Exporter{
			var rm AurSir4Go.AppMessage
			rm.Encode(rmsg,"JSON")
			c.appOutChannel <-AppMessage{exp,rm}
		}
	}
}


func (c core) dock(senderId string, dmsg AurSir4Go.AurSirDockMessage) {
	log.Println("Processing DOCK request from", senderId)
	c.storageAgent.Write(StorageCore.RegisterAppRequest{
		senderId,
		dmsg.AppName})
	var dm AurSir4Go.AppMessage
	dm.Encode(AurSir4Go.AurSirDockedMessage{}, "JSON")
	c.appOutChannel <- AppMessage{senderId, dm}
}

func (c core) leave(senderId string) {
	log.Println("Processing LEAVE request from", senderId)
	reply := c.storageAgent.Write(StorageCore.RemoveAppRequest{
		senderId})
	leave, ok := reply.(StorageCore.AppRemoved)
	if ok {

		for imp, appid := range leave.DisconnectedImports {
			var em AurSir4Go.AppMessage
			em.Encode(AurSir4Go.AurSirImportUpdatedMessage{imp, false}, "JSON")
			c.appOutChannel <- AppMessage{appid, em}
		}
	}

}

func (c core) addExport(senderId string, expMsg AurSir4Go.AurSirAddExportMessage) {
	log.Println("Processing ADD_EXPORT request from", senderId)
	reply := c.storageAgent.Write(StorageCore.AddExportRequest{senderId, expMsg.AppKey, expMsg.Tags})
	export, ok := reply.(StorageCore.ExportAdded)

	if ok {
		var em AurSir4Go.AppMessage
		em.Encode(AurSir4Go.AurSirExportAddedMessage{export.ExportId}, "JSON")
		c.appOutChannel <- AppMessage{senderId, em}
		for imp, appid := range export.ConnectedImports {
			var em AurSir4Go.AppMessage
			em.Encode(AurSir4Go.AurSirImportUpdatedMessage{imp, true}, "JSON")
			c.appOutChannel <- AppMessage{appid, em}
		}

		for _, r := range export.PendingJobs{
			var rm AurSir4Go.AppMessage

			rm.Encode(r,"JSON")


			c.appOutChannel <-AppMessage{senderId,rm}
		}
	}
}

func (c core) addImport(senderId string, expMsg AurSir4Go.AurSirAddImportMessage) {
	log.Println("Processing ADD_IMPORT request from", senderId)
	reply := c.storageAgent.Write(StorageCore.AddImportRequest{senderId, expMsg.AppKey, expMsg.Tags})
	imp, ok := reply.(StorageCore.ImportAdded)

	if ok {
		var em AurSir4Go.AppMessage
		em.Encode(AurSir4Go.AurSirImportAddedMessage{imp.ImportId, imp.Exported}, "JSON")
		c.appOutChannel <- AppMessage{senderId, em}
	}
}
