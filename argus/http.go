package argus

import (
	"net/http"
	"time"

	"github.com/boardware-cloud/common/code"
	argusModel "github.com/boardware-cloud/model/argus"
)

type HttpMonitor struct {
	entity argusModel.HttpMonitor
}

func (h *HttpMonitor) Sleep(a Argus) {
	lastRecord := a.Entity().LastRecord()
	if lastRecord == nil {
		return
	}
	time.Sleep(h.entity.Interval)
}

func (h *HttpMonitor) SetEntity(entity argusModel.Monitor) error {
	httpMonitor, ok := entity.(*argusModel.HttpMonitor)
	if !ok {
		return code.ErrConvert
	}
	h.entity = *httpMonitor
	return nil
}

func (h *HttpMonitor) Entity() argusModel.Monitor {
	return &h.entity
}

func (h *HttpMonitor) Check() Result {
	client := &http.Client{Timeout: time.Duration(h.entity.Timeout) * time.Second}
	req, _ := http.NewRequest(string(h.entity.HttpMethod), h.entity.Url, nil)
	// if m.Body != nil {
	// 	req, _ = http.NewRequest(string(*m.HttpMethod), m.Url, bytes.NewReader([]byte(*m.Body)))
	// }
	// if m.Headers != nil {
	// 	for _, header := range *m.Headers {
	// 		req.Header.Add(header.Left, header.Right)
	// 	}
	// }
	start := time.Now().UnixMilli()
	resp, err := client.Do(req)
	result := new(HttpCheckResult)
	result.responseTime = time.Duration(time.Now().UnixMilli() - start)
	if err != nil {
		if resp == nil {
			result.status = TIMEOUT
		} else {
			result.status = DOWN
		}
	} else {
		result.status = OK
		// record.StatusCode = fmt.Sprint(resp.StatusCode)
		// if checkAccepted(m.AcceptedStatusCodes, resp.StatusCode) {
		// 	result.status = constants.OK
		// } else {
		// 	record.Result = constants.DOWN
		// }
	}
	return result
}

type HttpCheckResult struct {
	status       ResultStatus
	responseTime time.Duration
}

func (r HttpCheckResult) Status() ResultStatus {
	return r.status
}

func (r HttpCheckResult) ResponseTime() time.Duration {
	return r.responseTime
}
