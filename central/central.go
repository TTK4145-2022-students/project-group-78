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

		state: MakeNetworkState(),
		stop:  make(chan bool),
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
			c.mergeNetworkState(s)
			c.StateOut <- c.state
		case <-c.stop:
			return
		}
	}
}

// Merge s into c.state
func (c *Central) mergeNetworkState(s NetworkState) {
	for id, es := range s {
		for event, time := range es {
			if time.After(c.state[id][event]) {
				c.state[id][event] = time
			}
		}
	}
}
