package argus

import (
	"time"

	"github.com/boardware-cloud/common/code"
	argusModel "github.com/boardware-cloud/model/argus"
	"github.com/go-ping/ping"
)

type PingMonitor struct {
	entity argusModel.PingMonitor
}

func (p PingMonitor) Sleep(a Argus) {
	lastRecord := a.Entity().LastRecord()
	if lastRecord == nil {
		return
	}
	time.Sleep(p.entity.Interval)
}

func (p *PingMonitor) SetEntity(entity argusModel.Monitor) error {
	pingMonitor, ok := entity.(*argusModel.PingMonitor)
	if !ok {
		return code.ErrConvert
	}
	p.entity = *pingMonitor
	return nil
}

func (p PingMonitor) Entity() argusModel.Monitor {
	return &p.entity
}

func (h *PingMonitor) Check() Result {
	pinger, err := ping.NewPinger(h.entity.Host)
	if err != nil {
		return &PingCheckResult{
			status: DOWN,
		}
	}
	pinger.Count = 3
	err = pinger.Run()
	if err != nil {
		return &PingCheckResult{
			status: DOWN,
		}
	}
	return &PingCheckResult{
		status: OK,
	}
}

type PingCheckResult struct {
	status ResultStatus
}

func (r PingCheckResult) Status() ResultStatus {
	return r.status
}

func (r PingCheckResult) ResponseTime() time.Duration {
	return 0
}

func (r *PingCheckResult) SetResponseTime() *PingCheckResult {
	return r
}
