package controllers

import (
	"net/http"
	"strconv"

	"flexoo-service-order/src/middlewares"
	"flexoo-service-order/src/models"
	"github.com/gin-gonic/gin"
)

type ServiceController struct{}

func NewServiceController() *ServiceController {
	return &ServiceController{}
}

func (sc *ServiceController) Explore(c *gin.Context) {
	client := middlewares.CurrentClient(c)
	categoryID, _ := strconv.ParseInt(c.Query("category"), 10, 64)
	search := c.Query("q")

	orders, _ := models.GetOrdersByClientID(client.ClientID)
	categories, _ := models.GetAllCategories()
	services, _ := models.GetAllServices(categoryID, search)
	pastWorks, _ := models.GetAllPastWorks()

	c.HTML(http.StatusOK, "explore.html", gin.H{
		"Client":             client,
		"Orders":             orders,
		"Categories":         categories,
		"Services":           services,
		"PastWorks":          pastWorks,
		"SelectedCategoryID": categoryID,
		"SearchQuery":        search,
		"OrderedID":          c.Query("ordered"),
	})
}

func (sc *ServiceController) ConfigureOrder(c *gin.Context) {
	client := middlewares.CurrentClient(c)
	serviceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	service, err := models.GetServiceByID(serviceID)
	if err != nil || service == nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "order.html", gin.H{
		"Client":  client,
		"Service": service,
	})
}

