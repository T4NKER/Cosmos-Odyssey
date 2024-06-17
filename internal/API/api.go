package API

import (
    "github.com/gin-gonic/gin"
    "Cosmos-Odyssey/internal/API/apis"
)

type MasterAPI struct {
    RoutesAPI      *apis.RouteAPI
    ReservationAPI *apis.ReservationAPI
}

func NewMasterAPI(routeAPI *apis.RouteAPI, reservationAPI *apis.ReservationAPI) *MasterAPI {
    return &MasterAPI{
        RoutesAPI:      routeAPI,
        ReservationAPI: reservationAPI,
    }
}

func (m *MasterAPI) RegisterRoutes(router *gin.Engine) {
    m.RoutesAPI.RegisterRoutes(router)
    m.ReservationAPI.RegisterRoutes(router)
    router.GET("/", m.homeHandler)
}

func (m *MasterAPI) homeHandler(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "Welcome to Cosomos Odyssey!",
    })
}