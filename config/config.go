package config

import (
	"fmt"
	"net"
	"time"
)

var BROADCAST_IP = net.ParseIP("127.255.255.255")

var DATAGRAM_PORT = 41784
var ACK_PORT = 41785
var HEARTBEAT_PORT = 41786

var TRANSMISSION_TIMEOUT = time.Second
var RETRANSMIT_INTERVAL = 10 * time.Millisecond

func LocalIp(id int) net.IP {
	return net.ParseIP(fmt.Sprintf("127.0.0.%v", id))
}
