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
	StateUpdate chan central.CentralState

	id   int
	conn *conn.Conn
	stop chan bool
}

func New(id int) *Distributor {
	d := &Distributor{
		StateUpdate: make(chan central.CentralState),

		id:   id,
		conn: conn.New(config.LocalIp(id), config.PORT),
		stop: make(chan bool),
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

func (d *Distributor) Send(c central.CentralState) {
	b := d.serialize(c)
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
			s, err := parse(b)
			if err != nil {
				d.log().WithField("packet", b).Warn(err)
			} else {
				d.log().Debug("received")
			}
			d.StateUpdate <- s

		case <-d.stop:
			return
		}
	}
}

func parse(b []byte) (s central.CentralState, err error) {
	err = gob.NewDecoder(bytes.NewBuffer(b)).Decode(&s)
	return
}

func (d *Distributor) serialize(s central.CentralState) []byte {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(s)
	if err != nil {
		d.log().Panic(err)
	}
	return buf.Bytes()
}
