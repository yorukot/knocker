package tasks

import (
	"strings"
)

const (
	TypeMonitorPingPattern = "monitor:ping:{region}"
)

// GetMonitorPingTypeForRegion returns the monitor ping task type for a specific region
func GetMonitorPingType(region string) string {
	return strings.Replace(TypeMonitorPingPattern, "{region}", region, 1)
}
