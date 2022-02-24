package order

import "github.com/TTK4145-2022-students/project-group-78/central"

type OrderType int

const (
	HallUp OrderType = iota
	HallDown
	Cab
)

type Order struct {
	Floor     int
	OrderType OrderType
}

func CalculateOrderLights(c central.CentralState) struct {
	Order
	bool
}                                                       {}
func CalculateTargetOrder(c central.CentralState) Order {}
