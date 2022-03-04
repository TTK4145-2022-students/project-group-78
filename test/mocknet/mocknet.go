package mocknet

import (
	"log"
	"math/rand"
	"net"
	"os/exec"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
	"github.com/tevino/abool"
)

func setDnat() {
	err := exec.Command("iptables", "-t", "nat", "-A", "OUTPUT", "-p", "udp", "-d", "255.255.255.255", "-j", "DNAT", "--to-destination", "127.255.255.255").Run()
	if err != nil {
		log.Panic(err)
	}
}

func unSetDnat() {
	err := exec.Command("iptables", "-t", "nat", "-D", "OUTPUT", "-p", "udp", "-d", "255.255.255.255", "-j", "DNAT", "--to-destination", "127.255.255.255").Run()
	if err != nil {
		log.Panic(err)
	}
}

type Mocknet struct {
	connectedPorts map[int]*abool.AtomicBool
	lossPercentage *int32
	stop           *abool.AtomicBool
}

func New(ports ...int) *Mocknet {
	setDnat()
	m := &Mocknet{
		connectedPorts: make(map[int]*abool.AtomicBool, len(ports)),
		lossPercentage: new(int32),
		stop:           abool.New(),
	}

	for _, port := range ports {
		m.connectedPorts[port] = abool.New()
		m.connectedPorts[port].Set()
		go m.run(port)
	}
	return m
}

// TODO: find way to stop goroutines

func (m *Mocknet) run(port int) {
	addr := &net.UDPAddr{IP: net.ParseIP("127.255.255.255"), Port: port}
	conn, err := net.ListenUDP("udp", addr)
	if err == nil {
		log.Print("mocknet: Listening on %v", conn.LocalAddr().String())
	} else {
		log.Panic(err)
	}

	msg := make([]byte, 1024)
	for m.stop.IsNotSet() {
		if _, err := conn.Read(msg); err != nil {
			if m.stop.IsSet() {
				break
			} else {
				log.Panic(err)
			}
		}

		if m.shouldLose() {
			continue
		}

		for port, connected := range m.connectedPorts {
			if connected.IsNotSet() {
				continue
			}
			if m.shouldLose() {
				continue
			}
			conn.WriteToUDP(msg, &net.UDPAddr{IP: nil, Port: port})
		}
	}
}

func (m *Mocknet) Stop() {
	m.stop.Set()
	unSetDnat()
}

func (m *Mocknet) SetLossPercentage(percentage int) {
	atomic.StoreInt32(m.lossPercentage, int32(percentage))
}

func (m *Mocknet) shouldLose() bool {
	return rand.Intn(100) <= int(atomic.LoadInt32(m.lossPercentage))
}
