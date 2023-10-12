package argus

import (
	argusModel "github.com/boardware-cloud/model/argus"
)

type Monitor interface {
	SetEntity(argusModel.Monitor) error
	Entity() argusModel.Monitor
	Sleep(Argus)
	Check() Result
}
