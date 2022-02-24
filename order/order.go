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

type OrderLight struct {
	Order Order
	Value bool
}

func CalculateOrderLights(c central.CentralState) []OrderLight {
	return []OrderLight{}
}

func CalculateTargetOrder(c central.CentralState) Order {
	return Order{}
}
