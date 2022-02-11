package conn

import (
	"net"

	log "github.com/sirupsen/logrus"
)

const MAX_PACKET_SIZE = 1024

type Conn struct {
	Send    chan []byte
	Receive chan []byte
	conn    *net.UDPConn
}

func New(localIp net.IP, localPort int, remoteIp net.IP, remotePort int) *Conn {
	localAddr := net.UDPAddr{IP: localIp, Port: localPort}
	remoteAddr := net.UDPAddr{IP: remoteIp, Port: remotePort}
	conn, err := net.DialUDP("udp", &localAddr, &remoteAddr)
	if err != nil {
		log.Panic(err)
	} else {
		log.Debugf("Dialed up %v from %v", remoteAddr.String(), localAddr.String())
	}

	c := &Conn{
		Send:   make(chan []byte, 100),
		Receive:  make(chan []byte, 100),
		conn: conn,
	}

	go c.sendForever()
	go c.receiveForever()

	return c
}

func (c *Conn) sendForever() {
	for {
		if len(c.Send) == cap(c.Send) {
			log.Panic("Send channel full")
		}
		c.send(<-c.Send)
	}
}

func (c *Conn) receiveForever() {
	for {
		if len(c.Receive) == cap(c.Receive) {
			log.Panic("Receive channel full")
		}
		c.Receive <- c.receive()
	}
}

func (c *Conn) send(packet []byte) {
	if len(packet) > MAX_PACKET_SIZE {
		log.Panicf("Packet size (%v) cannot exceed %v", len(packet), MAX_PACKET_SIZE)
	}

	n, err := c.conn.Write(packet)
	if err != nil {
		log.Panic(err)
	} else {
		log.Debugf("Sent %v bytes to %v", n, c.conn.RemoteAddr().String())
	}
}

func (c *Conn) receive() []byte {
	packet := make([]byte, MAX_PACKET_SIZE)
	n, addr, err := c.conn.ReadFromUDP(packet)
	if err != nil {
		log.Panic(err)
	} else {
		log.Debugf("Received %v bytes from %v", n, addr.String())
	}
	return packet[0:n]
}
