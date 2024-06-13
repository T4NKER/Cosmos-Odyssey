package API

import (
	"github.com/gin-gonic/gin"
	apis "internal/API/apis"
)

type MasterAPI struct {
	RoutesAPI      *RoutesAPI
	ReservationAPI *ReservationAPI
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

func (m *MasterAPI) homeHandler(router *gin.Engine) {

}