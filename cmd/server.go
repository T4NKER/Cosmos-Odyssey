package cmd

import (
	"Cosmos-Odyssey/internal/API"
	"Cosmos-Odyssey/internal/API/apis"
	"Cosmos-Odyssey/internal/services"
	"Cosmos-Odyssey/pkg/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start() {
	router := gin.Default()
	router.Use(cors.New(loadCorsConfig()))

	routeService := services.NewRouteService(database.Db)
	reservationService := services.NewReservationService()

	routeAPI := apis.NewRouteAPI(routeService)
	reservationAPI := apis.NewReservationAPI(reservationService)

	masterAPI := API.NewMasterAPI(routeAPI, reservationAPI)
	masterAPI.RegisterRoutes(router)

	router.Run()
}
