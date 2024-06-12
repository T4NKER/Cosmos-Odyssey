package apis

import (
	"github.com/gin-gonic/gin"
)

type MasterAPI struct {
	AuthAPI        *AuthAPI
	RoutesAPI      *RoutesAPI
	ReservationAPI *ReservationAPI
}

func NewMasterAPI(authAPI *AuthAPI, routesAPI *RoutesAPI, reservationAPI *ReservationAPI) *MasterAPI {
	return &MasterAPI{
		AuthAPI:        authAPI,
		RoutesAPI:      routesAPI,
		ReservationAPI: reservationAPI,
	}
}

func (m *MasterAPI) RegisterRoutes(router *gin.Engine) {
	m.AuthAPI.RegisterRoutes(router)
	m.RoutesAPI.RegisterRoutes(router)
	m.ReservationAPI.RegisterRoutes(router)
}