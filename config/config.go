package config

import "time"

const NumElevs = 3
const NumFloors = 4
const DoorOpenTime = 3 * time.Second
const TransmitInterval = 5 * time.Millisecond
const LightDelay = 200 * time.Millisecond
const ChanSize = 16
const OrderTimout = 200 * time.Millisecond
const ElevTimeout = DoorOpenTime + OrderTimout
