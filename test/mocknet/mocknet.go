package mocknet

import (
	"errors"
	"log"
	"math/rand"
	"net"
	"os/exec"
	"sync"
)

func setDnat() {
	err := exec.Command("sudo", "iptables", "-t", "nat", "-A", "OUTPUT", "-p", "udp", "-d", "255.255.255.255", "-j", "DNAT", "--to-destination", "127.255.255.255").Run()
	if err != nil {
		log.Panic(err)
	}
}

func unSetDnat() {
	err := exec.Command("sudo", "iptables", "-t", "nat", "-D", "OUTPUT", "-p", "udp", "-d", "255.255.255.255", "-j", "DNAT", "--to-destination", "127.255.255.255").Run()
	if err != nil {
		log.Panic(err)
	}
}

type Mocknet struct {
	conns          map[int]*net.UDPConn
	mutex          sync.Mutex
	lossPercentage int
}

func New(ports ...int) *Mocknet {
	setDnat()
	rand.Seed(0)
	m := &Mocknet{
		conns: make(map[int]*net.UDPConn),
		mutex: sync.Mutex{},
	}

	for _, port := range ports {
		m.Connect(port)
	}
	return m
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
	go m.run(conn)
}

func (m *Mocknet) run(conn *net.UDPConn) {
	for {
		msg := [1024]byte{}
		n, err := conn.Read(msg[:])
		if errors.Is(err, net.ErrClosed) {
			return
		} else if err != nil {
			log.Panic(err)
		}

		if m.shouldLose() {
			continue
		}

		for _, port := range m.getPorts() {
			if m.shouldLose() {
				continue
			}
			addr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: port}
			_, err := conn.WriteToUDP(msg[:n], addr)
			if errors.Is(err, net.ErrClosed) {
				return
			} else if err != nil {
				log.Panic(err)
			}
		}
	}
}

func (m *Mocknet) Stop() {
	for _, port := range m.getPorts() {
		m.Disconnect(port)
	}
	unSetDnat()
}

func (m *Mocknet) Disconnect(port int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.conns[port].Close()
	delete(m.conns, port)
}

func (m *Mocknet) SetLossPercentage(percentage int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.lossPercentage = percentage
}

func (m *Mocknet) shouldLose() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return rand.Intn(100) <= m.lossPercentage
}

func (m *Mocknet) getPorts() []int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	ports := make([]int, len(m.conns))
	i := 0
	for port := range m.conns {
		ports[i] = port
		i++
	}
	return ports
}
