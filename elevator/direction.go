package elevator

import "github.com/TTK4145-2022-students/driver-go-group-78/elevio"

type Direction int

const (
	Up Direction = iota
	Down
)

func (d Direction) toMd() elevio.MotorDirection {
	return map[Direction]elevio.MotorDirection{Up: elevio.MD_Up, Down: elevio.MD_Down}[d]
}

func (d Direction) opposite() Direction {
	return map[Direction]Direction{Up: Down, Down: Up}[d]
}

func (d Direction) toBt() elevio.ButtonType {
	return map[Direction]elevio.ButtonType{Up: elevio.BT_HallUp, Down: elevio.BT_HallDown}[d]
}

func (d Direction) ToString() string {
	return map[Direction]string{Up: "up", Down: "down"}[d]
}
