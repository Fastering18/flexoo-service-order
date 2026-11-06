package models

import (
	"fmt"
	"math/rand"
	"time"

	"flexoo-service-order/src/config"
)

type ServiceOrder struct {
	OrderID           string    `json:"order_id"`
	ClientID          int64     `json:"client_id"`
	ServiceID         int64     `json:"service_id"`
	ServiceName       string    `json:"service_name"`
	ProjectName       string    `json:"project_name"`
	Scale             string    `json:"scale"`
	Comments          string    `json:"comments"`
	SRSFileURL        string    `json:"srs_file_url"`
	Status            string    `json:"status"`
	EstimatedDelivery time.Time `json:"estimated_delivery"`
	CreatedAt         time.Time `json:"created_at"`
}

func (o ServiceOrder) FormattedDeliveryDate() string {
	if o.EstimatedDelivery.IsZero() {
		return "TBD"
	}
	return o.EstimatedDelivery.Format("02 January 2006")
}

func GetOrdersByClientID(clientID int64) ([]ServiceOrder, error) {
	if config.DB == nil {
		storeMu.RLock()
		defer storeMu.RUnlock()
		var list []ServiceOrder
		for i := len(fallbackOrders) - 1; i >= 0; i-- {
			if fallbackOrders[i].ClientID == clientID {
				list = append(list, fallbackOrders[i])
			}
		}
		return list, nil
	}

	rows, err := config.DB.Query(
		`SELECT o.order_id, o.client_id, o.service_id, COALESCE(s.name, ''), o.project_name, o.scale, COALESCE(o.comments, ''), COALESCE(o.srs_file_url, ''), o.status, o.estimated_delivery, o.created_at
		 FROM service_orders o
		 LEFT JOIN services s ON o.service_id = s.service_id
		 WHERE o.client_id = $1
		 ORDER BY o.created_at DESC, o.order_id DESC`,
		clientID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []ServiceOrder
	for rows.Next() {
		var o ServiceOrder
		if err := rows.Scan(&o.OrderID, &o.ClientID, &o.ServiceID, &o.ServiceName, &o.ProjectName, &o.Scale, &o.Comments, &o.SRSFileURL, &o.Status, &o.EstimatedDelivery, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func CreateServiceOrder(order *ServiceOrder) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	order.OrderID = fmt.Sprintf("ORDER-%d", 1000+r.Intn(9000))
	order.Status = "In Progress"
	order.CreatedAt = time.Now()

	days := 14
	if order.Scale == "MVP" {
		days = 7
	} else if order.Scale == "Enterprise (custom)" {
		days = 30
	}
	order.EstimatedDelivery = time.Now().AddDate(0, 0, days)

	if config.DB == nil {
		if svc, _ := GetServiceByID(order.ServiceID); svc != nil {
			order.ServiceName = svc.Name
		}
		storeMu.Lock()
		fallbackOrders = append(fallbackOrders, *order)
		storeMu.Unlock()
		return nil
	}

	_, err := config.DB.Exec(
		`INSERT INTO service_orders (order_id, client_id, service_id, project_name, scale, comments, srs_file_url, status, estimated_delivery, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		order.OrderID, order.ClientID, order.ServiceID, order.ProjectName, order.Scale, order.Comments, order.SRSFileURL, order.Status, order.EstimatedDelivery, order.CreatedAt,
	)
	return err
}

func UpdateOrderStatus(orderID string, clientID int64, status string) error {
	if config.DB == nil {
		storeMu.Lock()
		defer storeMu.Unlock()
		for i := range fallbackOrders {
			if fallbackOrders[i].OrderID == orderID && fallbackOrders[i].ClientID == clientID {
				fallbackOrders[i].Status = status
				return nil
			}
		}
		return fmt.Errorf("order not found")
	}

	_, err := config.DB.Exec(`UPDATE service_orders SET status = $1 WHERE order_id = $2 AND client_id = $3`, status, orderID, clientID)
	return err
}

