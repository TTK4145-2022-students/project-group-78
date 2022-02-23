package central

type Central struct {
	StateIn  chan NetworkState
	StateOut chan NetworkState

	state NetworkState
	stop  chan bool
}

func New() *Central {
	c := &Central{
		StateIn:  make(chan NetworkState),
		StateOut: make(chan NetworkState),

		state: make(NetworkState),
		stop: make(chan bool),
	}

	go c.run()

	return c
}

func (c *Central) Stop() {
	c.stop <- true
}

func (c *Central) run() {
	for {
		select {
		case s := <-c.StateIn:
			mergeNetworkState(c.state, s)
			c.StateOut <- c.state
		case <-c.stop:
			return
		}
	}
}

// Merge s2 into s
func mergeNetworkState(s1 NetworkState, s2 NetworkState) {
	for id, es := range s2 {
		mergeElevatorState(s1[id], es)
	}
}

// Merge es2 into es
func mergeElevatorState(es ElevatorState, es2 ElevatorState) {
	for event, time := range es2 {
		if time.After(es[event]) {
			es[event] = time
		}
	}
}
