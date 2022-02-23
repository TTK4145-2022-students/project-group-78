package distributor

import "github.com/TTK4145-2022-students/project-group-78/events"

type Distributor struct {
}

func New(id int) *Distributor {
	return &Distributor{}
}

func (d *Distributor) RegisterComponent(in chan events.Event, out chan events.Event, []Type{}) {

}

func (d *Distributor) Start() {
	for {}
}

func (d *Distributor) Stop()
