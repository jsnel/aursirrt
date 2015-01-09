package dockzmq

import (
	"log"
	"github.com/joernweissenborn/aursirrt/core"
	"github.com/joernweissenborn/aursir4go"
	zmq "github.com/pebbe/zmq4"


	"strconv"
	"net"
	"time"
	"dock"
	"github.com/joernweissenborn/aursir4go/messages"
)

type DockerZmq struct {
	agent dock.DockAgent
	skt *zmq.Socket
}

func (dzmq DockerZmq) Launch(agent dock.DockAgent) (err error) {

	mprint("Launching")
	dzmq.agent = agent
	dzmq.skt, err = zmq.NewSocket(zmq.ROUTER)

	if err != nil {
		mprint("Failed to start ZMQ Dock")
		return
	}

	dzmq.Bind("tcp://*:5555")

	go dzmq.listen()
	return
}

func (dzmq *DockerZmq) listen() {

	mprint("ZMQAppDocker listening")

	for {

		msg, _ := dzmq.skt.RecvMessage(0)

		if len(msg) > 3 {

			senderId := msg[0]

			msgtype, err := strconv.ParseInt(msg[1], 10, 64)
			log.Println("ZMQAppDocker got message from", senderId)
			codec := msg[2]
			if err == nil {

				switch msgtype {

				case messages.DOCK:
					encmsg := []byte(msg[3])
					port, err := strconv.ParseInt(msg[1], 10, 64)
					if err == nil {
						conn := NewConnection(port)
						dzmq.agent.InitDocking(senderId, codec, encmsg, conn)
					}
				default:
					encmsg := []byte(msg[3])
					dzmq.agent.ProcessMsg(senderId,msgtype,codec,encmsg)

				}
			}
		}
	}
}


func (dzmq *DockerZmq) openConnection(id, port string, c chan core.AppMessage){

	skt, _ := zmq.NewSocket(zmq.DEALER)
	defer skt.Close()
	skt.SetIdentity("AURSIR_RT")

	skt.SetSndtimeo(1000)

	log.Println("ZMQAppDocker opening channel to port", port)

	skt.Connect("tcp://localhost:" + port)



	for msg := range c {


		appmsg := msg.AppMsg
		_,err := skt.SendMessage(
			[]string{strconv.FormatInt(appmsg.MsgType,10),appmsg.MsgCodec,string(appmsg.Msg)},0)


		if err != nil {
			log.Println("ZMQAppDocker Error on zqm port %d, closing:",port,err)
			dzmq.closeConnection(id)
		}

	}
	log.Println("DOCKERZMQ", "Closing connection to port",port)


}

func (dzmq *DockerZmq)updPingListener() {
	var buf [1024]byte
	addr, err := net.ResolveUDPAddr("udp", ":5556")
	if err != nil {
		log.Fatal("DOCKERZMQ",err)
	}
	sock, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("DOCKERZMQ",err)
	}
	log.Println("DOCKERZMQ", "Startet to listen for UDP ping on port 5556")
	for {
		rlen, _, err := sock.ReadFromUDP(buf[:])
		if err != nil {
			log.Fatal("DOCKERZMQ",err)
		}
		appid := string(buf[:rlen])
		//log.Println("PING",appid)
		go dzmq.CheckInAppPing(appid)
	}
}

func (dzmq *DockerZmq) AddAppPing(AppId string) {
	dzmq.appLastPing[AppId] = time.Now()
}

func (dzmq *DockerZmq) CheckInAppPing(AppId string) {
	if _,f:= dzmq.appLastPing[AppId];f {
		dzmq.appLastPing[AppId] = time.Now()
	}
}

func (dzmq *DockerZmq) RemoveAppPing(AppId string) {
	delete(dzmq.appLastPing,AppId)
}

func (dzmq *DockerZmq) CheckAppLiveliness() {

	t := time.NewTimer(10 * time.Second)
	for _ = range t.C {
		for id, lastCheckIn := range dzmq.appLastPing {
			if time.Since(lastCheckIn) > 10*time.Second {
				dzmq.closeConnection(id)
			}
		}
		t.Reset(10*time.Second)
	}
}

func (dzmq *DockerZmq) closeConnection(id string){
	l := []byte("{}")
	dzmq.RemoveAppPing(id)
	dzmq.msgChan <- core.AppMessage{id,aursir4go.AppMessage{aursir4go.LEAVE,"JSON",l}}
	dzmq.regChan <- ungisterDockedApp{id}
}

func mprint(msg string){
	log.Println("DOCKERZMQ", msg)
}

func printDebug(msg string){
	log.Println("DEBUG","DOCKERZMQ", msg)
}
