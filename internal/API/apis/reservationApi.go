package apis

import (
	"Cosmos-Odyssey/internal/models"
	"Cosmos-Odyssey/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReservationAPI struct {
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
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "There was an error while parsing the reservation information"})
		log.Println(err)
		return
	}

	err := r.ReservationService.ValidateReservation(reservation)
	if err != nil {
			if err.Error() == "invalid reservation, there's a new pricelist already, refresh and try again" {
				c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
				return
			}
			c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "There was an error validating the reservation"})
			return
	}

	reservationSuccess, err := r.ReservationService.MakeReservation(reservation)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "There was an error creating the reservation"})
		return
	}

	c.HTML(http.StatusOK, "reservation.html", gin.H{"error": nil, "reservationSuccess": reservationSuccess})
}
