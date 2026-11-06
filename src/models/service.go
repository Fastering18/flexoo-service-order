package models

import (
	"fmt"
	"strings"

	"flexoo-service-order/src/config"
)

type Service struct {
	ServiceID    int64                `json:"service_id"`
	CategoryID   int64                `json:"category_id"`
	CategoryName string               `json:"category_name"`
	Name         string               `json:"name"`
	Description  string               `json:"description"`
	PageRange    string               `json:"page_range"`
	BasePrice    float64              `json:"base_price"`
	Deliverables []ServiceDeliverable `json:"deliverables,omitempty"`
	PastWorks    []PastWork           `json:"past_works,omitempty"`
}

func (s Service) FormattedPrice() string {
	if s.BasePrice == float64(int64(s.BasePrice)) {
		return fmt.Sprintf("$%d", int64(s.BasePrice))
	}
	return fmt.Sprintf("$%.2f", s.BasePrice)
}

func matchesFilter(s Service, categoryID int64, search string) bool {
	if categoryID > 0 && s.CategoryID != categoryID {
		return false
	}
	if search == "" {
		return true
	}
	q := strings.ToLower(search)
	return strings.Contains(strings.ToLower(s.Name), q) ||
		strings.Contains(strings.ToLower(s.Description), q) ||
		strings.Contains(strings.ToLower(s.CategoryName), q)
}

func GetAllServices(categoryID int64, search string) ([]Service, error) {
	search = strings.TrimSpace(search)

	if config.DB == nil {
		storeMu.RLock()
		defer storeMu.RUnlock()
		var list []Service
		for _, s := range fallbackServices {
			if matchesFilter(s, categoryID, search) {
				list = append(list, s)
			}
		}
		return list, nil
	}

	rows, err := config.DB.Query(
		`SELECT s.service_id, COALESCE(s.category_id, 0), COALESCE(c.category_name, ''), s.name, COALESCE(s.description, ''), COALESCE(s.page_range, ''), s.base_price
		 FROM services s
		 LEFT JOIN categories c ON s.category_id = c.category_id
		 ORDER BY s.service_id ASC`,
	)
	if err != nil {
		return fallbackServices, err
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var s Service
		if err := rows.Scan(&s.ServiceID, &s.CategoryID, &s.CategoryName, &s.Name, &s.Description, &s.PageRange, &s.BasePrice); err != nil {
			return nil, err
		}
		if matchesFilter(s, categoryID, search) {
			services = append(services, s)
		}
	}
	return services, nil
}

func GetServiceByID(serviceID int64) (*Service, error) {
	services, err := GetAllServices(0, "")
	if err != nil {
		return nil, err
	}
	for _, s := range services {
		if s.ServiceID == serviceID {
			s.Deliverables, _ = GetDeliverablesByServiceID(serviceID)
			s.PastWorks, _ = GetPastWorksByServiceID(serviceID)
			return &s, nil
		}
	}
	return nil, fmt.Errorf("service not found")
}

