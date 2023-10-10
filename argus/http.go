package argus

import (
	"time"

	argusModel "github.com/boardware-cloud/model/argus"
)

type HttpMonitor struct {
	Monitor argusModel.HttpMonitor
}

func (h HttpMonitor) Interval() time.Duration {
	return time.Duration(h.Monitor.Interval) * time.Second
}

func (h HttpMonitor) Sleep() {
	time.Sleep(h.Interval() * time.Second)
}
