package config

import "time"

const (
	NumElevs          = 3
	NumFloors         = 4
	NumOrderTypes     = 3
	NumHallOrderTypes = 2
	DoorOpenTime      = 3 * time.Second
	TransmitInterval  = 15 * time.Millisecond
	LightDelay        = 300 * time.Millisecond
	ChanSize          = 16
	OrderTimout       = 500 * time.Millisecond
	ElevTimeout       = DoorOpenTime + OrderTimout
	BcastPort         = 56986
)
