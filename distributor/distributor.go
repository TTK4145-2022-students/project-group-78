package distributor

import (
	"bytes"
	"encoding/gob"

	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/conn"
	"github.com/TTK4145-2022-students/project-group-78/utils"
	"github.com/sirupsen/logrus"
)

var Logger = utils.NewLogger("distributor", "id")

type Distributor struct {
	id       byte
	conn     *conn.Conn
	stateOut chan central.NetworkState
	stop     chan bool
}

func New(id byte, stateOut chan central.NetworkState) *Distributor {
	d := &Distributor{
		id:       id,
		conn:     conn.New(config.LocalIp(id), config.PORT),
		stateOut: make(chan central.NetworkState),
		stop:     make(chan bool),
	}

	go d.run()
	d.log().Info("started")

	return d
}

func (d *Distributor) Stop() {
	d.stop <- true
	d.conn.Close()
	d.log().Info("stopped")
}

func (d *Distributor) Send(s central.NetworkState) {
	b := d.serialize(s)
	d.conn.SendTo(b, config.BROADCAST_IP, config.PORT)
	d.log().Debug("sent")
}

func (d *Distributor) log() *logrus.Entry {
	return Logger.WithField("id", d.id)
}

func (d *Distributor) run() {
	for {
		select {
		case b := <-d.conn.Receive:
			ns, err := parse(b)
			if err != nil {
				d.log().WithField("packet", b).Warn(err)
			} else {
				d.log().Debug("received")
			}
			d.stateOut <- ns

		case <-d.stop:
			return
		}
	}
}

func parse(b []byte) (s central.NetworkState, err error) {
	err = gob.NewDecoder(bytes.NewBuffer(b)).Decode(&s)
	return
}

func (d *Distributor) serialize(s central.NetworkState) []byte {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(s)
	if err != nil {
		d.log().Panic(err)
	}
	return buf.Bytes()
}
