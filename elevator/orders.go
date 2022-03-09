package elevator

import "github.com/TTK4145-2022-students/project-group-78/config"

type Orders [config.NUM_FLOORS][3]bool

func (o Orders) Above(floor int) bool {
	for f := floor + 1; f < len(o); f++ {
		for bt := 0; bt < 3; bt++ {
			if o[f][bt] {
				return true
			}
		}
	}
	return false
}

func (o Orders) Below(floor int) bool {
	for f := 0; f < floor; f++ {
		for bt := 0; bt < 3; bt++ {
			if o[f][bt] {
				return true
			}
		}
	}
	return false
}

func (o Orders) Here(f int) bool {
	for bt := 0; bt < 3; bt++ {
		if o[f][bt] {
			return true
		}
	}
	return false
}