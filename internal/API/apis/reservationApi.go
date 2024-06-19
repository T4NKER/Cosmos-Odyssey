package apis

import (
	"Cosmos-Odyssey/internal/services"
	"Cosmos-Odyssey/internal/models"
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
	if err := c.Bind(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservationSuccess, err := r.ReservationService.MakeReservation(reservation)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.HTML(http.StatusOK, "reservation.html", reservationSuccess)
}