package argus

import (
	"time"

	argusModel "github.com/boardware-cloud/model/argus"
)

type Monitor interface {
	SetEntity(argusModel.Monitor) error
	Entity() argusModel.Monitor
	Sleep()
	Interval() time.Duration
	Check()
	Alive() bool
}
