package config

import "time"

const (
	NumElevs          = 3
	NumFloors         = 4
	NumOrderTypes     = 3
	NumHallOrderTypes = 2
	DoorOpenTime      = 3 * time.Second
	TransmitInterval  = 10 * time.Millisecond
	LightDelay        = 33 * TransmitInterval
	ChanSize          = 16
	OrderTimout       = 500 * time.Millisecond
	ElevTimeout       = DoorOpenTime + OrderTimout
	BcastPort         = 56986
)
