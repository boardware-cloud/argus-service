package argus

import (
	"time"

	argusModel "github.com/boardware-cloud/model/argus"
)

type PingMonitor struct {
	Monitor argusModel.PingMonitor
}

func (p PingMonitor) Interval() time.Duration {
	return time.Duration(p.Monitor.Interval) * time.Second
}

func (p PingMonitor) Sleep() {
	time.Sleep(p.Interval() * time.Second)
}
