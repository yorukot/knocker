package statuspage

import (
	"fmt"
)

// validateStatusPagePayload enforces sort order uniqueness and group references.
func validateStatusPagePayload(req statusPageUpsertRequest) error {
	groupOrders := make(map[int]struct{})
	groupIDs := make(map[int64]struct{})
	fatherOrders := make(map[int]struct{})

	for _, g := range req.Groups {
		if _, exists := groupOrders[g.SortOrder]; exists {
			return fmt.Errorf("duplicate group sort_order %d", g.SortOrder)
		}
		groupOrders[g.SortOrder] = struct{}{}
		if _, exists := fatherOrders[g.SortOrder]; exists {
			return fmt.Errorf("duplicate top-level sort_order %d between groups and ungrouped monitors", g.SortOrder)
		}
		fatherOrders[g.SortOrder] = struct{}{}
		if g.ID != nil {
			groupIDs[*g.ID] = struct{}{}
		}
	}

	monitorOrders := make(map[string]map[int]struct{}) // key: groupID|nil
	for _, m := range req.Monitors {
		groupKey := "nil"
		groupLabel := "ungrouped monitors"
		if m.GroupID != nil {
			if _, ok := groupIDs[*m.GroupID]; !ok {
				return fmt.Errorf("monitor group_id %d not present in groups list", *m.GroupID)
			}
			groupKey = fmt.Sprintf("%d", *m.GroupID)
			groupLabel = fmt.Sprintf("group %s", groupKey)
		} else {
			if _, exists := fatherOrders[m.SortOrder]; exists {
				return fmt.Errorf("duplicate top-level sort_order %d between groups and ungrouped monitors", m.SortOrder)
			}
			fatherOrders[m.SortOrder] = struct{}{}
		}

		if _, ok := monitorOrders[groupKey]; !ok {
			monitorOrders[groupKey] = make(map[int]struct{})
		}
		if _, exists := monitorOrders[groupKey][m.SortOrder]; exists {
			return fmt.Errorf("duplicate monitor sort_order %d within %s", m.SortOrder, groupLabel)
		}
		monitorOrders[groupKey][m.SortOrder] = struct{}{}
	}

	if err := validateConsecutiveOrders("top-level", fatherOrders); err != nil {
		return err
	}
	for groupKey, orders := range monitorOrders {
		if groupKey != "nil" {
			label := fmt.Sprintf("group %s monitors", groupKey)
			if err := validateConsecutiveOrders(label, orders); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateConsecutiveOrders(label string, orders map[int]struct{}) error {
	if len(orders) == 0 {
		return nil
	}
	for order := range orders {
		if order < 1 {
			return fmt.Errorf("%s sort_order must start at 1 and be consecutive", label)
		}
	}
	for i := 1; i <= len(orders); i++ {
		if _, ok := orders[i]; !ok {
			return fmt.Errorf("%s sort_order must be consecutive starting at 1", label)
		}
	}
	return nil
}
