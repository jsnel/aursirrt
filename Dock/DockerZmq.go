package dock

import (
	"log"
	"github.com/joernweissenborn/aursirrt/core"
	"github.com/joernweissenborn/aursir4go"
	zmq "github.com/pebbe/zmq4"


	"strconv"
)

type DockerZmq struct {
	msgChan chan core.AppMessage
	regChan chan registerDockedApp
}

func (dzmq DockerZmq) Launch(mc chan core.AppMessage, rc chan registerDockedApp) {

	log.Println("DockerZMQ Launching")

	dzmq.msgChan = mc
	dzmq.regChan = rc

	incoming, err := zmq.NewSocket(zmq.ROUTER)

	if err != nil {
		log.Fatal("Failed to start ZMQ Dock")
		panic(err)
	}

	incoming.Bind("tcp://*:5555")

	go dzmq.listen(incoming)

	log.Println("ZMQAppDocker launched")

}

func (dzmq DockerZmq) listen(skt *zmq.Socket) {

	log.Println("ZMQAppDocker listening")

	for {

		msg, _ := skt.RecvMessage(0)

		if len(msg) >3{

		senderId := msg[0]


		log.Println("ZMQAppDocker got message from", senderId)
		msgtype,err := strconv.ParseInt(msg[1],10,64)
		if err ==nil{

			if msgtype == aursir4go.DOCK{
				c := make(chan core.AppMessage )

				dzmq.regChan <- registerDockedApp{senderId, c}
				go dzmq.openConnection(msg[4],c)
			}
			encmsg := []byte(msg[3])
			dzmq.msgChan <- core.AppMessage{senderId,aursir4go.AppMessage{msgtype,msg[2],&encmsg}}
		}
	}
	}
}

func (dzmq DockerZmq) openConnection(port string, c chan core.AppMessage){

	skt, _ := zmq.NewSocket(zmq.DEALER)
	defer skt.Close()
	skt.SetIdentity("AURSIR_RT")

	skt.SetSndtimeo(1000)

	log.Println("ZMQAppDocker opening channel to port", port)

	skt.Connect("tcp://localhost:" + port)



	for msg := range c {


		appmsg := msg.AppMsg
		_,err := skt.SendMessage(
			[]string{strconv.FormatInt(appmsg.MsgType,10),appmsg.MsgCodec,string(*appmsg.Msg)},0)


		if err != nil {
			log.Println("ZMQAppDocker Error on zqm port %d, closing:",port,err)
		}

	}

}
