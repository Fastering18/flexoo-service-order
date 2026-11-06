package handler

import (
	"net/http"

	"flexoo-service-order/src/routes"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	router := routes.SetupRouter()
	router.ServeHTTP(w, r)
}
