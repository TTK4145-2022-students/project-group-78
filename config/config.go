package config

import "time"

const (
	NumElevs          = 3
	NumFloors         = 4
	NumOrderTypes     = 3
	NumHallOrderTypes = 2

	DoorOpenTime     = 3 * time.Second
	TransmitInterval = 15 * time.Millisecond
	OrderTimout      = 300 * time.Millisecond
	LightDelay       = OrderTimout
	ElevTimeout      = DoorOpenTime + OrderTimout

	ChanSize  = 16
	BcastPort = 58735
)
