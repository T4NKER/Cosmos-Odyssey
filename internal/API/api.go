package API

import (
	"Cosmos-Odyssey/internal/API/apis"

	"github.com/gin-gonic/gin"
)

type MasterAPI struct {
	RoutesAPI      *apis.RouteAPI
	ReservationAPI *apis.ReservationAPI
	HomeAPI        *apis.HomeAPI
}

func NewMasterAPI(routeAPI *apis.RouteAPI, reservationAPI *apis.ReservationAPI, homeAPI *apis.HomeAPI) *MasterAPI {
	return &MasterAPI{
		RoutesAPI:      routeAPI,
		ReservationAPI: reservationAPI,
		HomeAPI:        homeAPI,
	}
}

func (m *MasterAPI) RegisterRoutes(router *gin.Engine) {
	m.RoutesAPI.RegisterRoutes(router)
	m.ReservationAPI.RegisterRoutes(router)
	m.HomeAPI.RegisterRoutes(router)
}
