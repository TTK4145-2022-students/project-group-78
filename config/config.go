package config

import (
	"fmt"
	"net"
)

var BROADCAST_IP = net.ParseIP("127.255.255.255")
const PORT = 41784

const NUM_ELEVATORS = 3
const NUM_FLOORS = 4

func LocalIp(id byte) net.IP {
	return net.ParseIP(fmt.Sprintf("127.0.0.%v", id))
}
