package models

import "flexoo-service-order/src/config"

type PastWork struct {
	WorkID    int64  `json:"work_id"`
	ServiceID int64  `json:"service_id"`
	Title     string `json:"title"`
	ImageURL  string `json:"image_url"`
}

func GetAllPastWorks() ([]PastWork, error) {
	if config.DB == nil {
		storeMu.RLock()
		defer storeMu.RUnlock()
		return append([]PastWork(nil), fallbackPastWorks...), nil
	}

	rows, err := config.DB.Query(`SELECT work_id, service_id, title, image_url FROM past_works ORDER BY work_id ASC`)
	if err != nil {
		return fallbackPastWorks, err
	}
	defer rows.Close()

	var works []PastWork
	for rows.Next() {
		var w PastWork
		if err := rows.Scan(&w.WorkID, &w.ServiceID, &w.Title, &w.ImageURL); err != nil {
			return nil, err
		}
		works = append(works, w)
	}
	return works, nil
}

func GetPastWorksByServiceID(serviceID int64) ([]PastWork, error) {
	all, err := GetAllPastWorks()
	if err != nil {
		return nil, err
	}
	var list []PastWork
	for _, w := range all {
		if w.ServiceID == serviceID {
			list = append(list, w)
		}
	}
	if len(list) < 2 && len(all) >= 2 {
		return all[:2], nil
	}
	return list, nil
}

