package central

import "time"

type Event interface{}

type NetworkState map[byte]ElevatorState

type ElevatorState map[Event]time.Time
