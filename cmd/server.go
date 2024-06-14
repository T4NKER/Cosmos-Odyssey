package server

import (
	"github.com/T4NKER/Cosmos-Odyssey/internal/API"
	"github.com/T4NKER/Cosmos-Odyssey/internal/API/apis"
	"github.com/T4NKER/Cosmos-Odyssey/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

)

func Start() {
	router := gin.Default()
	router.Use(cors.New(loadCorsConfig()))

	routeService := services.NewRouteService()
	reservationService := services.NewReservationService()

	routeAPI := apis.NewRouteAPI(routeService)
	reservationAPI := apis.NewReservationAPI(reservationService)

	masterAPI := API.NewMasterAPI(routeAPI, reservationAPI)
	masterAPI.RegisterRoutes(router)

	router.Run()
}


