package models

import "flexoo-service-order/src/config"

type Category struct {
	CategoryID   int64  `json:"category_id"`
	CategoryName string `json:"category_name"`
}

func GetAllCategories() ([]Category, error) {
	if config.DB == nil {
		storeMu.RLock()
		defer storeMu.RUnlock()
		return append([]Category(nil), fallbackCategories...), nil
	}

	rows, err := config.DB.Query(`SELECT category_id, category_name FROM categories ORDER BY category_id ASC`)
	if err != nil {
		return fallbackCategories, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.CategoryID, &c.CategoryName); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

