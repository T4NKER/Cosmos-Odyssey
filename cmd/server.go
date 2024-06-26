package cmd

import (
	"Cosmos-Odyssey/internal/API"
	"Cosmos-Odyssey/internal/API/apis"
	"Cosmos-Odyssey/internal/services"
	"Cosmos-Odyssey/internal/services/external"
	"Cosmos-Odyssey/pkg/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start() {
	router := gin.Default()
	router.LoadHTMLGlob("./frontend/*.html")
	router.Use(cors.New(loadCorsConfig()))

	pricelistService := external.NewExternalPricelistService()
	routeService := services.NewRouteService(database.Db, pricelistService)
	reservationService := services.NewReservationService(database.Db, pricelistService)
	homeService := services.NewHomeService(database.Db)

	routeAPI := apis.NewRouteAPI(routeService)
	reservationAPI := apis.NewReservationAPI(reservationService)
	homeAPI := apis.NewHomeAPI(homeService)

	masterAPI := API.NewMasterAPI(routeAPI, reservationAPI, homeAPI)
	masterAPI.RegisterRoutes(router)

	router.Run()
}
