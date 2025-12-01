package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/yorukot/knocker/models"
)

// MonitorPingPayload represents the payload for a monitor ping task
type MonitorPingPayload struct {
	Monitor models.Monitor `json:"monitor"`
	Region  string         `json:"region"`
}

func NewMonitorPing(monitor models.Monitor, region string) (*asynq.Task, error) {
	payload := MonitorPingPayload{
		Monitor: monitor,
		Region:  region,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(GetMonitorPingType(region), payloadBytes), nil
}
