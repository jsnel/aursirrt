/*
	Package Core is responsible for processing incoming messages
*/
package core

import (
	"github.com/joernweissenborn/AurSir4Go"
	"github.com/joernweissenborn/aursirrt/storagecore"
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
	storageAgent storagecore.StorageCoreAgent
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

			case  AurSir4Go.AurSirUpdateExportMessage:
				go c.updateExport(AppMessage.SenderUUID,aursirMessage)

			case AurSir4Go.AurSirAddImportMessage:
				go c.addImport(AppMessage.SenderUUID, aursirMessage)
			case  AurSir4Go.AurSirUpdateImportMessage:
				go c.updateImport(AppMessage.SenderUUID,aursirMessage)

			case AurSir4Go.AurSirRequest:
				go c.request(AppMessage.SenderUUID,aursirMessage)

			case AurSir4Go.AurSirResult:
				go c.result(AppMessage.SenderUUID,aursirMessage)
			case AurSir4Go.AurSirListenMessage:
				go c.listen(AppMessage.SenderUUID,aursirMessage)
			case AurSir4Go.AurSirStopListenMessage:
				go c.stopListen(AppMessage.SenderUUID,aursirMessage)

			case AurSir4Go.AurSirCallChain:
				go c.callChain(AppMessage.SenderUUID,aursirMessage)

			default:
				log.Println("unknown message")
			}
		}
	}
}

func (c core) callChain(senderId string,ccmsg AurSir4Go.AurSirCallChain) {
	log.Println("Processing CALL_CHAIN request from", senderId)


	oak := ccmsg.OriginRequest.AppKeyName
	ofn := ccmsg.OriginRequest.FunctionName
	ok := true
	insane := []int64{}
	for i,call := range ccmsg.CallChain {
		 if !c.chainChecker(oak,ofn,call.AppKeyName,call.FunctionName,call.ArgumentMap){
			 ok = false
			 insane = append(insane,int64(i))
		 }
	}
	log.Println("CALLCHAIN is vailid:",ok, insane)

	var ccm AurSir4Go.AppMessage
	ccm.Encode(AurSir4Go.AurSirCallChainAddedMessage{ok,insane},"JSON")
	c.appOutChannel <-AppMessage{senderId,ccm}

	if ok {
		exports, ok := (c.storageAgent.Write(storagecore.AddCallChainRequest{senderId,ccmsg})).([]string)
		if ok {
			for _, export := range exports {
				var reqmsg AurSir4Go.AppMessage
				reqmsg.Encode(ccmsg.OriginRequest,"JSON")
				c.appOutChannel <-AppMessage{export,reqmsg}
			}
		}
	}
}

func (c core) chainChecker(orgAppKey, orgFun, tarAppKey, tarFun string, paramap map[string]string) bool {
	oak, f := (c.storageAgent.Read(storagecore.GetAppKey{orgAppKey})).(AurSir4Go.AppKey)
	log.Println(oak)
	if !f {
		return false
	}
	tak, f := c.storageAgent.Read(storagecore.GetAppKey{tarAppKey}).(AurSir4Go.AppKey)
	if !f  {
		return false
	}
	log.Println(tak)

	var ofn AurSir4Go.Function
	f = false
	for _,fkt := range oak.Functions{
		if fkt.Name == orgFun {
			ofn = fkt
			f = true
		}
	}

	if !f {return false}
	tmp := map[string]int{}
	for input,output := range paramap {
		f =false
		for _,out := range ofn.Output {
			if out.Name == output {
				f=true
				tmp[input] = out.Type
		}
			if !f {return false}
		}
	}
	var tfn AurSir4Go.Function
	f = false
	for _,fkt := range tak.Functions{
		if fkt.Name == tarFun {
			tfn = fkt
			f = true
		}
	}

	if !f {return false}
	for _,in :=range tfn.Input {
		t,f := tmp[in.Name]
		if !f || t != in.Type {
			return false
		}
	}
	return true
}


func (c core) updateImport(senderId string,uimsg AurSir4Go.AurSirUpdateImportMessage) {
	log.Println("Processing UPDATE_IMPORT request from", senderId)
	reply := c.storageAgent.Write(storagecore.UpdateImportRequest{uimsg})
	imp, ok := reply.(storagecore.ImportAdded)

	if ok {
		var em AurSir4Go.AppMessage
		em.Encode(AurSir4Go.AurSirImportUpdatedMessage{imp.ImportId, imp.Exported}, "JSON")
		c.appOutChannel <- AppMessage{senderId, em}
	}
}

func (c core) updateExport(senderId string,uemsg AurSir4Go.AurSirUpdateExportMessage){
	log.Println("Processing UPDATE_EXPORT request from", senderId)
	reply := c.storageAgent.Write(storagecore.UpdateExportRequest{uemsg})
	export, ok := reply.(storagecore.ExportAdded)
	log.Println(export)
	if ok {
		for imp, appid := range export.ConnectedImports {
			var em AurSir4Go.AppMessage
			em.Encode(AurSir4Go.AurSirImportUpdatedMessage{imp, true}, "JSON")
			c.appOutChannel <- AppMessage{appid, em}
			log.Println("sending to",appid)
		}


		for imp, appid := range export.DisconnectedImports {
			var em AurSir4Go.AppMessage
			em.Encode(AurSir4Go.AurSirImportUpdatedMessage{imp, false}, "JSON")
			c.appOutChannel <- AppMessage{appid, em}
		}
		for _, r := range export.PendingJobs{
			var rm AurSir4Go.AppMessage
			rm.Encode(r,"JSON")
			c.appOutChannel <-AppMessage{senderId,rm}
		}
	}
}

func (c core) listen(senderId string, lmsg AurSir4Go.AurSirListenMessage){
	log.Println("Processing LISTEN request from", senderId)

	c.storageAgent.Write(storagecore.ListenRequest{senderId,lmsg.FunctionName,lmsg.ImportId})
}

func (c core) stopListen(senderId string, slmsg AurSir4Go.AurSirStopListenMessage){
	log.Println("Processing STOP_LISTEN request from", senderId)

}

func (c core) result(senderId string, rmsg AurSir4Go.AurSirResult) {
	log.Println("Processing RESULT request from", senderId)
	reply := c.storageAgent.Write(storagecore.AddResRequest{senderId,rmsg})
	resReg, ok := reply.(storagecore.ResRegistered)

	if ok && len(resReg.Importer)!=0{
		for _, imp := range resReg.Importer{
			var rm AurSir4Go.AppMessage
			rm.Encode(rmsg,"JSON")
			c.appOutChannel <-AppMessage{imp,rm}
		}
	}

	if ok && resReg.IsChainCall{
		c.createChainCall(senderId,rmsg,resReg.ChainCall,resReg.ChainCallImportId)
	}
}

func (c core) createChainCall(senderId string, prevResult AurSir4Go.AurSirResult,cc AurSir4Go.ChainCall,ccImportId string){

	codec := AurSir4Go.GetCodec(prevResult.Codec)
	if codec==nil {
		return
	}
	var tmp interface {}

	codec.Decode(&prevResult.Result,&tmp)
	log.Println(string(prevResult.Result),tmp)
	resultParameter := tmp.(map[string]interface {})
	requestParameter := map[string]interface {}{}

	for target,origin := range cc.ArgumentMap{
		requestParameter[target]=resultParameter[origin]
	}
	req , err:=codec.Encode(requestParameter)
	if err != nil {
		return
	}
	request := AurSir4Go.AurSirRequest{
		cc.AppKeyName,
		cc.FunctionName,
		cc.CallType,
		cc.Tags,
		cc.ChainCallId,
		ccImportId,
		prevResult.Codec,
		*req}

	reqReg, ok := c.storageAgent.Write(storagecore.AddReqRequest{senderId,request}).(storagecore.ReqRegistered)
	log.Println(ok,reqReg)
	if ok{
		for _, exp := range reqReg.Exporter {
			var rm AurSir4Go.AppMessage
			rm.Encode(request, "JSON")
			c.appOutChannel <-AppMessage{exp, rm}
		}
	}
}

func (c core) request(senderId string, rmsg AurSir4Go.AurSirRequest) {
	log.Println("Processing REQUEST request from", senderId)
	reply := c.storageAgent.Write(storagecore.AddReqRequest{senderId,rmsg})
	reqreg, ok := reply.(storagecore.ReqRegistered)
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
	c.storageAgent.Write(storagecore.RegisterAppRequest{
		senderId,
		dmsg.AppName})
	var dm AurSir4Go.AppMessage
	dm.Encode(AurSir4Go.AurSirDockedMessage{}, "JSON")
	c.appOutChannel <- AppMessage{senderId, dm}
}

func (c core) leave(senderId string) {
	log.Println("Processing LEAVE request from", senderId)
	reply := c.storageAgent.Write(storagecore.RemoveAppRequest{
		senderId})
	leave, ok := reply.(storagecore.AppRemoved)
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
	reply := c.storageAgent.Write(storagecore.AddExportRequest{senderId, expMsg.AppKey, expMsg.Tags})
	export, ok := reply.(storagecore.ExportAdded)

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
	reply := c.storageAgent.Write(storagecore.AddImportRequest{senderId, expMsg.AppKey, expMsg.Tags})
	imp, ok := reply.(storagecore.ImportAdded)

	if ok {
		var em AurSir4Go.AppMessage
		em.Encode(AurSir4Go.AurSirImportAddedMessage{imp.ImportId, imp.Exported, expMsg.AppKey.ApplicationKeyName,expMsg.Tags}, "JSON")
		c.appOutChannel <- AppMessage{senderId, em}
	}
}
