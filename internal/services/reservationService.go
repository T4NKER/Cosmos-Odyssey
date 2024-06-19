package services

import (
	"Cosmos-Odyssey/internal/models"
	"database/sql"
	"github.com/google/uuid"
	"log"
)

type ReservationService struct {
	DB *sql.DB
}

func NewReservationService(database *sql.DB) *ReservationService {
	return &ReservationService{
		DB: database,
	}
}

func (r *ReservationService) MakeReservation(reservation models.Reservation) (models.Reservation, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		log.Println("Failed to begin transaction: ", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to error: ", err)
		} else {
			err = tx.Commit()
			if err != nil {
				log.Println("Failed to commit transaction: ", err)
			}
		}
	}()

	newID := uuid.New().String()
	reservation.ID = newID

	newRoute := r.createNewRoute(reservation.Route)

	_, err = tx.Exec("INSERT INTO reservations (id, pricelist_id, first_name, last_name, route_id, total_quoted_price, total_quoted_travel_time, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		reservation.ID, reservation.PricelistID, reservation.Firstname, reservation.Lastname, newRoute, reservation.TotalQuotedPrice, reservation.TotalQuotedTravelTime, reservation.TransportationCompanyNames)
	if err != nil {
		log.Println("Failed to insert reservation: ", err)
		return reservation, err
	}

	return reservation, nil
}

func (r *ReservationService) createNewRoute(route []string) string {
	newRoute := ""
	for _, v := range route {
		newRoute += ", " + v
	}
	return newRoute
}
