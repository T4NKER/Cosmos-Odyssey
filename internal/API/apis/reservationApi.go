package apis

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"cosmos-odyssey/models"
)

type ReservationAPI struct{
	ReservationService *ReservationService
}

func NewReservationAPI(reservationService *ReservationService) *ReservationAPI {
	return &ReservationAPI{ReservationService: reservationService}
}

func (r *ReservationAPI) RegisterRoutes(router *gin.Engine) {
	router.POST("/reservation", r.MakeReservation)
}

func (r *ReservationAPI) MakeReservation(c *gin.Context) {
	var reservation models.Reservation
	if err := c.BindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Implement logic to save reservation

	c.JSON(http.StatusOK, gin.H{"success": true})
}