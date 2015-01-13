package dockzmq

import (
	"log"
	zmq "github.com/pebbe/zmq4"
	"strconv"
	"aursirrt/src/dock"
	"github.com/joernweissenborn/aursir4go/messages"
	"net"
	"time"
	"encoding/json"
	"strings"
	"fmt"
)

type DockerZmq struct {
	agent dock.DockAgent
	skt *zmq.Socket
	id string
	ip string
	appLastPing map[string]time.Time
	homeport int64
	p2p bool
	killPing chan struct {}
}

func (dzmq DockerZmq) Launch(agent dock.DockAgent, id string) (err error) {

	mprint(fmt.Sprint("Launching on port",dzmq.homeport))
	dzmq.agent = agent
	dzmq.id = id
	dzmq.appLastPing = map[string]time.Time{}
	dzmq.skt, err = zmq.NewSocket(zmq.ROUTER)

	if err != nil {
		mprint("Failed to start ZMQ Dock")
		return
	}

	dzmq.skt.Bind("tcp://*:"+strconv.FormatInt(dzmq.homeport,10))

	go dzmq.listen()
	if dzmq.p2p {
		dzmq.launchUdp()
	}
	return
}

func (dzmq *DockerZmq) SetP2P(p2p bool){
	mprint(fmt.Sprint("P2P active: "))
	dzmq.p2p = p2p
}
func (dzmq *DockerZmq) SetPort(homeport int64){
	mprint(fmt.Sprint("Incoming port is: ",homeport))
	dzmq.homeport = homeport
}
func (dzmq *DockerZmq) SetIp(ip string){
	mprint(fmt.Sprint("Broadcastinterfaceadress is: ",ip))
	dzmq.ip = ip
}

func (dzmq *DockerZmq) listen() {

	mprint("Listening for incoming")

	for {

		msg, _ := dzmq.skt.RecvMessage(0)

		if len(msg) > 5 {

			senderId := msg[0]

			msgtype, err := strconv.ParseInt(msg[1], 10, 64)
			log.Println("ZMQAppDocker got message from", msg)
			codec := msg[2]
			if err == nil {

				switch msgtype {

				case messages.DOCK:
					encmsg := []byte(msg[3])
					IP := "localhost"
					if len(msg) > 6{
					 IP = msg[5]
					}
					printDebug(IP)
					port, err := strconv.ParseInt(msg[4], 10, 64)
					if err == nil {
						conn := NewConnection(dzmq.homeport, port, dzmq.ip,IP, dzmq.id)
						dzmq.addAppPing(senderId)
						dzmq.agent.InitDocking(senderId, codec, encmsg, &conn)
					}
				default:
					encmsg := []byte(msg[3])
					dzmq.agent.ProcessMsg(senderId,msgtype,codec,encmsg)

				}
			}
		}
	}
}

func (dzmq *DockerZmq)launchUdp() {

	go dzmq.updPingListener()
	go dzmq.checkAppLiveliness()
	kill := make(chan struct {})
	go pingUdp(dzmq.id, dzmq.ip, dzmq.homeport,kill)

}
func (dzmq *DockerZmq)updPingListener() {
	var buf [1024]byte



	ifaceaddresses,_ := net.InterfaceAddrs()
	var Interface net.Interface
	for i,iface := range ifaceaddresses {

		addr := strings.Split(iface.String(),"/")[0]
		mprint(fmt.Sprint("found networkinterface ",addr))
		if addr == dzmq.ip {
			ifaces,_ := net.Interfaces()
			Interface = ifaces[i]
		}
	}

	mcip, err := net.ResolveUDPAddr("udp", "224.0.0.251:5556")

	sock, err := net.ListenMulticastUDP("udp4", &Interface, mcip)
	if err != nil {
		log.Fatal("DOCKERZMQ",err)
	}
	log.Println("DOCKERZMQ", "Startet to listen for UDP ping on port 5556")
	for {
		rlen, Ip, err := sock.ReadFromUDP(buf[:])
		if err != nil {
			log.Fatal("DOCKERZMQ",err)
		}
		beaconstring := strings.Split(string(buf[:rlen]),":")
		appid := beaconstring[0]
		port := ""
		if len(beaconstring) >1{
			port = beaconstring[1]
		}
		ip := strings.Split(Ip.String(),":")[0]

		go dzmq.checkInAppPing(appid, ip,port)
	}
}

func (dzmq *DockerZmq) addAppPing(AppId string) {
	dzmq.appLastPing[AppId] = time.Now()
}

func (dzmq *DockerZmq) checkInAppPing(AppId string, ip string, Port string) {
	if _,f:= dzmq.appLastPing[AppId];f {
		dzmq.appLastPing[AppId] = time.Now()
	} else if AppId!= dzmq.id && Port != ""{
		port,_ := strconv.ParseInt(Port,10,64)
		conn := NewConnection( dzmq.homeport, port, dzmq.ip, ip, dzmq.id)
		err := conn.Init()
		if err != nil {
			return
		}
		defer conn.Close()
		m,_ := json.Marshal(messages.DockMessage{"runtime@"+ip,[]string{"JSON"},true})
		conn.Send(messages.DOCK,"JSON",m)
	}
}

func (dzmq *DockerZmq) removeAppPing(AppId string) {
	delete(dzmq.appLastPing,AppId)
}

func (dzmq *DockerZmq) checkAppLiveliness() {

	t := time.NewTimer(10 * time.Second)
	for _ = range t.C {
		for id, lastCheckIn := range dzmq.appLastPing {
			if time.Since(lastCheckIn) > 10*time.Second {
				mprint("Apptimeout: "+id)
				dzmq.closeConnection(id)
			}
		}
		t.Reset(10*time.Second)
	}
}

func (dzmq *DockerZmq) closeConnection(id string){

	dzmq.removeAppPing(id)
	dzmq.agent.ProcessMsg(id,messages.LEAVE,"JSON",[]byte{})
}



func pingUdp(UUID, ip string, port int64, kill chan struct {}) {

	var pingtime = 8 * time.Second

	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:0",ip))
	if err != nil {
		log.Fatal("DOCKERZMQ", err)
	}
	serverAddr, err := net.ResolveUDPAddr("udp", "224.0.0.251:5556")

	if err != nil {
		log.Fatal("DOCKERZMQ", err)
	}
	con, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		log.Fatal("DOCKERZMQ", err)
	}
	t := time.NewTimer(pingtime)
	mprint(fmt.Sprintf("Beginning UDP Broadcast on %s:%d",ip,port))
	for {
		select {
		case <-kill:
			mprint("Stopping UDP")
			return

		case <-t.C:
			con.Write([]byte(fmt.Sprintf("%s:%d", UUID, port)))
			t.Reset(pingtime)
		}
	}

}

func mprint(msg string){
	log.Println("DOCKERZMQ", msg)
}

func printDebug(msg string){
	log.Println("DEBUG","DOCKERZMQ", msg)
}
