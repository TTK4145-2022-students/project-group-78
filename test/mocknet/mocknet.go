package mocknet

import (
	"errors"
	"log"
	"math/rand"
	"net"
	"time"
)

type Mocknet struct {
	LossPercentage int

	conns   map[int]*net.UDPConn
	closed  bool
	msgC    chan []byte
	closedC chan int
}

func New(ports ...int) *Mocknet {
	rand.Seed(0)
	m := &Mocknet{
		conns:   make(map[int]*net.UDPConn, len(ports)),
		msgC:    make(chan []byte),
		closedC: make(chan int),
	}
	for _, port := range ports {
		m.Connect(port)
	}
	go m.run()
	return m
}

func (m *Mocknet) run() {
	for {
		select {
		case closed := <-m.closedC:
			delete(m.conns, closed)

		case msg := <-m.msgC:
			if shouldLose(m.LossPercentage) {
				continue
			}
			for port, conn := range m.conns {
				if !shouldLose(m.LossPercentage) {
					send(conn, port, msg)
				}
			}

		case <-time.After(time.Second):
			if m.closed && len(m.conns) == 0 {
				log.Print("mocknet: Closed")
				return
			}
		}
	}
}

func (m *Mocknet) Connect(port int) {
	addr := &net.UDPAddr{IP: net.ParseIP("127.255.255.255"), Port: port}
	conn, err := net.ListenUDP("udp", addr)
	if err == nil {
		log.Printf("mocknet: Listening on %v", conn.LocalAddr().String())
	} else {
		log.Panic(err)
	}
	m.conns[port] = conn
	go receiver(conn, port, m.msgC, m.closedC)
}

func (m *Mocknet) Disconnect(port int) {
	conn, ok := m.conns[port]
	if !ok {
		log.Panicf("mocknet: Can not disconnect %v because it is not connected", port)
	}
	if err := conn.Close(); err != nil {
		log.Panic(err)
	}
}

func (m *Mocknet) Close() {
	for _, conn := range m.conns {
		err := conn.Close()
		if errors.Is(err, net.ErrClosed) {
			continue
		} else if err != nil {
			log.Panic(err)
		}
	}
	m.closed = true
}

func receiver(conn *net.UDPConn, port int, msgC chan []byte, closedC chan int) {
	for {
		msg := [1024]byte{}
		n, err := conn.Read(msg[:])
		switch {
		case errors.Is(err, net.ErrClosed):
			log.Print(err)
			closedC <- port
			return

		case err != nil:
			log.Panic(err)

		default:
			msgC <- msg[:n]
		}
	}
}

func shouldLose(lossPercentage int) bool {
	return rand.Intn(100) < lossPercentage
}

func send(conn *net.UDPConn, port int, msg []byte) {
	addr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: port}
	if _, err := conn.WriteToUDP(msg, addr); err != nil {
		log.Panic(err)
	}
}
