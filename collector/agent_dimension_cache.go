package collector

import (
	"sync"
	"time"
)

type AgentDimensionRefresher struct {
}

const (
	TooManyRequestsErrorCode = 429
)

var GetAgentDimensionRefresher = NewAgentDimensionRefresherGetter()

func NewAgentDimensionRefresherGetter() func() *AgentDimensionRefresher {
	var once sync.Once
	var refresher *AgentDimensionRefresher
	return func() *AgentDimensionRefresher {
		once.Do(func() {
			refresher = &AgentDimensionRefresher{}
		})
		return refresher
	}
}

func (r *AgentDimensionRefresher) RefreshAgentDimensionWithInterval() {
	go func() {
		r.GetAgentDimension()
		t := time.NewTicker(time.Duration(10) * time.Minute)
		for range t.C {
			r.GetAgentDimension()
		}
	}()
}

func (r *AgentDimensionRefresher) GetAgentDimension() {
	syncInstanceAgentDimensions(&ecsInfo)
	syncInstanceAgentDimensions(&bmsInfo)
}

func syncInstanceAgentDimensions(srvsInfo *serversInfo) {
	var instanceIDs []string
	srvsInfo.Lock()
	if srvsInfo.LabelInfo != nil {
		for instanceID := range srvsInfo.LabelInfo {
			instanceIDs = append(instanceIDs, instanceID)
		}
	}
	srvsInfo.Unlock()
	// 发生限流时可重试，重试前等待时间递增，最多60s
	waitSecondArrayForRetry := []int{1, 5, 10, 30, 60, 0}
	for _, instanceID := range instanceIDs {
		for _, waitSecond := range waitSecondArrayForRetry {
			err := loadAgentDimensions(instanceID)
			if err != nil && isErrorTypeForTooManyRequests(err) {
				time.Sleep(time.Duration(waitSecond) * time.Second)
				continue
			}
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
}
