// services control center
package ftunnel

import (
	"log"
	"net"
)

type Core struct {
	Nodes          []Node
	Services       []Service
	BinaryUrl      string
	BinaryCheckSum []byte
	listener       net.Listener
}

func (co *Core) Start() {
	// determine which node is this node
	myip := ip()
	for _, nd := range co.Nodes {
		if myip == nd.Ip {
			nd.Identity = _nodeId
			co.StartListen(nd.Port)
		} else {
			// check identity of other node
			go nd.Connect()
		}
	}

	// start all services
	for _, s := range co.Services {
		s.co = co
		s.Start()
	}
}

func (co *Core) Stop() {
	// close all nodes connections
	for _, n := range co.Nodes {
		n.Close()
	}

	co.listener.Close()
	// close all services
	for _, s := range co.Services {
		//if is tcplistener
		s.Stop()
	}
}

func (co *Core) StartListen(port string) {
	var err error
	co.listener, err = net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("E(core.StartListen):", err)
		return
	}

	for {
		conn, err := co.listener.Accept()
		if err != nil {
			// handle error
			log.Println("N(core.Accept):", err)
			continue
		}
		tr := NewTransporter(conn)
		go tr.ServConnection()
	}
}

func (co *Core) Send(b []byte, s *Service, connid uint64) {
	// TODO:
	// if this node is inbound or outbound
	// if this is outbound send one to dst
	// 
}
