package argus

import (
	"net/http"
	"strconv"
	"strings"
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
	for _, header := range h.entity.Headers {
		req.Header.Add(header.Left, header.Right)
	}
	tries := int64(0)
	result := new(HttpCheckResult)

	for tries <= h.entity.Retries {
		tries++
		start := time.Now()
		resp, err := client.Do(req)
		result.SetResponseTime(time.Since(start))
		if err != nil {
			if resp == nil {
				result.status = TIMEOUT
			} else {
				result.status = DOWN
			}
		} else {
			if checkAccepted(h.entity.AcceptedStatusCodes, resp.StatusCode) {
				result.status = OK
			} else {
				result.status = DOWN
			}
		}
	}
	return result
}

func checkAccepted(acceptedStatusCodes []string, statusCode int) bool {
	if len(acceptedStatusCodes) != 0 {
		for _, code := range acceptedStatusCodes {
			if checkAcceptedStatusCode(code, statusCode) {
				return true
			}
		}
		return false
	}
	if statusCode >= 200 && statusCode < 300 {
		return true
	}
	return false
}

func checkAcceptedStatusCode(acceptedStatusCode string, statusCode int) bool {
	codes := strings.Split(acceptedStatusCode, "-")
	if len(codes) == 1 {
		return codes[0] == strconv.Itoa(statusCode)
	}
	left, err := strconv.Atoi(codes[0])
	if err != nil {
		return false
	}
	right, err := strconv.Atoi(codes[len(codes)-1])
	if err != nil {
		return false
	}
	if left > right {
		temp := left
		left = right
		right = temp
	}
	if statusCode >= left && statusCode <= right {
		return true
	}
	return false
}

type HttpCheckResult struct {
	status       ResultStatus
	responseTime time.Duration
}

func (h HttpCheckResult) Status() ResultStatus {
	return h.status
}

func (h HttpCheckResult) ResponseTime() time.Duration {
	return h.responseTime
}

func (h *HttpCheckResult) SetResponseTime(r time.Duration) *HttpCheckResult {
	h.responseTime = r
	return h
}
