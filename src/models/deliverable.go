package models

import "flexoo-service-order/src/config"

type ServiceDeliverable struct {
	DeliverableID int64  `json:"deliverable_id"`
	ServiceID     int64  `json:"service_id"`
	Label         string `json:"label"`
}

func GetDeliverablesByServiceID(serviceID int64) ([]ServiceDeliverable, error) {
	if config.DB == nil {
		storeMu.RLock()
		defer storeMu.RUnlock()
		var list []ServiceDeliverable
		for _, d := range fallbackDeliverables {
			if d.ServiceID == serviceID {
				list = append(list, d)
			}
		}
		return list, nil
	}

	rows, err := config.DB.Query(
		`SELECT deliverable_id, service_id, label FROM service_deliverables WHERE service_id = $1 ORDER BY deliverable_id ASC`,
		serviceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ServiceDeliverable
	for rows.Next() {
		var d ServiceDeliverable
		if err := rows.Scan(&d.DeliverableID, &d.ServiceID, &d.Label); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, nil
}

