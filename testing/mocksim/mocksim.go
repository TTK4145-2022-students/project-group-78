package mocksim

import (
	"log"
	"net"
	"sync/atomic"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
	"github.com/tevino/abool"
)

func toByte(a bool) byte {
	if a {
		return byte(1)
	} else {
		return byte(0)
	}
}

func toBool(a byte) bool {
	return a != 0
}

type OrderLight struct {
	Floor int
	Type  elevio.ButtonType
	Value bool
}

func SetFloor(f int) {
	atomic.StoreInt32(floor, int32(f))
}

var OrderButtons [config.NUM_FLOORS][3]*abool.AtomicBool
var AtFloor *abool.AtomicBool
var Stop *abool.AtomicBool
var Obstruction *abool.AtomicBool
var floor *int32

func Sim(port int, reloadConfigC chan bool, motorDirectionC chan elevio.MotorDirection, orderLightC chan OrderLight, floorIndicatorC chan int, doorLightC chan bool, stopLightC chan bool) {
	for f := 0; f < len(OrderButtons); f++ {
		for i := 0; i < len(OrderButtons[f]); i++ {
			OrderButtons[f][i] = abool.New()
		}
	}
	AtFloor = abool.New()
	Stop = abool.New()
	Obstruction = abool.New()
	floor = new(int32)

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: nil, Port: port})
	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := listener.AcceptTCP()

		for {
			msg := make([]byte, 4)
			_, err = conn.Read(msg)
			if err != nil {
				log.Print(err)
				break
			}

			switch msg[0] {
			case 0:
				reloadConfigC <- true

			case 1:
				motorDirectionC <- elevio.MotorDirection(msg[1])

			case 2:
				orderLightC <- OrderLight{Floor: int(msg[2]), Type: elevio.ButtonType(msg[1]), Value: toBool(msg[3])}

			case 3:
				floorIndicatorC <- int(msg[1])

			case 4:
				doorLightC <- toBool(msg[1])

			case 5:
				stopLightC <- toBool(msg[1])

			case 6:
				value := OrderButtons[int(msg[2])][int(msg[1])].IsSet()
				err = send(conn, msg[0], toByte(value))

			case 7:
				err = send(conn, msg[0], toByte(AtFloor.IsSet()), byte(atomic.LoadInt32(floor)))

			case 8:
				err = send(conn, msg[0], toByte(Stop.IsSet()))

			case 9:
				err = send(conn, msg[0], toByte(Obstruction.IsSet()))
			}
			if err != nil {
				log.Print(err)
				break
			}
		}
	}
}

func send(conn *net.TCPConn, msg ...byte) error {
	_, err := conn.Write(msg)
	return err
}
