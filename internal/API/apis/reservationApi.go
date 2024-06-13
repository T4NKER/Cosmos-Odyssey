package apis

import (
	"Cosomos-Odyssey/internal/services"
	"cosmos-odyssey/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReservationAPI struct{
	ReservationService *services.ReservationService
}

func NewReservationAPI(reservationService *services.ReservationService) *ReservationAPI {
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

	err := r.ReservationService.MakeReservation()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	// Implement logic to save reservation

	c.JSON(http.StatusOK, gin.H{"success": true})
}