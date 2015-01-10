package dockzmq

import (
	"strconv"
	zmq "github.com/pebbe/zmq4"
	"fmt"
)


type ConnectionZmq struct {
	port int64
	skt *zmq.Socket
}

func NewConnection(port int64) ConnectionZmq{
	return ConnectionZmq{port, nil}
}

func (cz *ConnectionZmq) Init() (err error) {
	cz.skt, _ = zmq.NewSocket(zmq.DEALER)
	cz.skt.SetIdentity("AURSIR_RT")

	cz.skt.SetSndtimeo(1000)
	port := strconv.FormatInt(cz.port,10)
	printDebug("ZMQAppDocker opening channel to port "+port)

	err = cz.skt.Connect("tcp://localhost:" + port)

	return

}

func (cz ConnectionZmq) Send(msgtype int64, codec string,msg []byte) (err error){
	_,err = cz.skt.SendMessage(
		[]string{strconv.FormatInt(msgtype,10),codec,string(msg)},0)

	if err != nil {
		mprint(fmt.Sprintf("Error on zqm port %d, closing:",cz.port,err))
	}
	return
}
