package config

import "time"

const NUM_ELEVS = 3
const NUM_FLOORS = 4
const DOOR_OPEN_TIME = 3 * time.Second
const TRANSMIT_INTERVAL = 5 * time.Millisecond
const LIGHT_DELAY = 167 * time.Millisecond
const CHAN_SIZE = 16
const ORDER_TIMEOUT = 200 * time.Millisecond
const ELEV_TIMEOUT = 3200 * time.Millisecond

//PEERS: 15 ms / 500 ms
// 500/ 3 = 167
