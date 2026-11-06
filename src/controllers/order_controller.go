package controllers

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"flexoo-service-order/src/middlewares"
	"flexoo-service-order/src/models"
	"github.com/gin-gonic/gin"
)

type OrderController struct{}

func NewOrderController() *OrderController {
	return &OrderController{}
}

func (oc *OrderController) PlaceOrder(c *gin.Context) {
	client := middlewares.CurrentClient(c)

	serviceID, err := strconv.ParseInt(c.PostForm("service_id"), 10, 64)
	if err != nil || serviceID <= 0 {
		serviceID = 1
	}

	scale := strings.TrimSpace(c.PostForm("scale"))
	if scale == "" {
		scale = "Standard"
	}

	projectName := strings.TrimSpace(c.PostForm("project_name"))
	if projectName == "" {
		projectName = "FLEXOO PROJECT"
	}

	srsFileName := strings.TrimSpace(c.PostForm("srs_file_name"))
	if fileHeader, err := c.FormFile("srs_file"); err == nil && fileHeader != nil {
		srsFileName = fileHeader.Filename
	}
	if srsFileName == "" {
		srsFileName = "demo_srs_document.pdf"
	}

	order := &models.ServiceOrder{
		ClientID:    client.ClientID,
		ServiceID:   serviceID,
		ProjectName: projectName,
		Scale:       scale,
		Comments:    strings.TrimSpace(c.PostForm("comments")),
		SRSFileURL:  srsFileName,
	}

	_ = models.CreateServiceOrder(order)
	c.Redirect(http.StatusSeeOther, "/?ordered="+url.QueryEscape(order.OrderID))
}

func (oc *OrderController) UpdateStatus(c *gin.Context) {
	client := middlewares.CurrentClient(c)
	orderID := c.Param("id")
	status := strings.TrimSpace(c.PostForm("status"))
	if status == "" {
		status = "In Progress"
	}

	_ = models.UpdateOrderStatus(orderID, client.ClientID, status)
	c.Redirect(http.StatusSeeOther, "/")
}

