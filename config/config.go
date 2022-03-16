package config

import "time"

const NUM_ELEVS = 1
const NUM_FLOORS = 4
const DOOR_OPEN_TIME = 3 * time.Second
const TRANSMIT_INTERVAL = 5 * time.Millisecond
const CHAN_SIZE = 16
const ORDER_TIMEOUT = time.Second
const ELEV_TIMEOUT = 3 * time.Second
