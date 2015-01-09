package dock

import "processor/processors"

import "processor"

import (
)



type DockAgent struct {
	procchan chan processor.Processor
}

func (da DockAgent) ProcessMsg(appid string,msgtype int64, codec string, msg []byte){
	var p processors.ParseMessageProccesor
	p.AppId = appid
	p.Codec = codec
	p.Msg = msg
	p.Type = msgtype
	go da.launchProcess(p)
}

func (da DockAgent) InitDocking(appid string, codec string, msg []byte, connection Connection){
	var p processors.DockProcessor
	p.AppId = appid
	p.Codec = codec
	p.DockMessage = msg
	p.Connection = connection
	go da.launchProcess(p)

}

func (da DockAgent) launchProcess(p processor.Processor){
	da.procchan <- p
}

