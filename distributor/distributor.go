package distributor

import (
	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/conn"
)

type Distributor struct {
	id byte
	conn *conn.Conn
	stateOut chan central.NetworkState
	stop  chan bool
}

func New(id byte, stateOut chan central.NetworkState) *Distributor {
	c := &Distributor{
		conn: conn.New(config.LocalIp())
		stateOut: make(chan central.NetworkState),

		stop:  make(chan bool),
	}

	go c.run()

	return c
}

func (c *Distributor) Stop() {
	c.stop <- true
}

func (d *Distributor) Send(ns central.NetworkState) {

}

func (c *Distributor) run() {
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
