package statuspage

import "fmt"

func normalizeStatusPageUpsert(req statusPageUpsertRequest) (statusPageUpsertRequest, error) {
	if len(req.Elements) == 0 {
		return req, nil
	}

	groups := make([]statusPageGroupInput, 0, len(req.Elements))
	monitors := make([]statusPageMonitorInput, 0)
	nextTempID := int64(-1)

	for _, element := range req.Elements {
		if element.Monitor {
			if element.MonitorID == nil {
				return req, fmt.Errorf("monitor_id is required for monitor elements")
			}
			monitors = append(monitors, statusPageMonitorInput{
				ID:        nil,
				MonitorID: *element.MonitorID,
				GroupID:   nil,
				Name:      element.Name,
				Type:      element.Type,
				SortOrder: element.SortOrder,
			})
			continue
		}

		groupID := element.ID
		if groupID == nil {
			tempID := nextTempID
			nextTempID--
			groupID = &tempID
		}

		groups = append(groups, statusPageGroupInput{
			ID:        groupID,
			Name:      element.Name,
			Type:      element.Type,
			SortOrder: element.SortOrder,
		})

		for _, monitor := range element.Monitors {
			monitor.GroupID = groupID
			monitors = append(monitors, monitor)
		}
	}

	req.Groups = groups
	req.Monitors = monitors
	req.Elements = nil

	return req, nil
}
