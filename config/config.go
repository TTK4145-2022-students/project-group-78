package config

import "net"

var BROADCAST_IP = net.ParseIP("127.255.255.255")

var DATAGRAM_PORT = 41784
var ACK_PORT = 41785
