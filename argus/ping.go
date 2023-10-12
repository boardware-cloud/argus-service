package argus

import (
	"time"

	"github.com/boardware-cloud/common/code"
	argusModel "github.com/boardware-cloud/model/argus"
)

type PingMonitor struct {
	entity argusModel.PingMonitor
}

func (p PingMonitor) Sleep(a Argus) {
	// TODO
	time.Sleep(5 * time.Second)
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
	// TODO
	return &PingCheckResult{}
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
