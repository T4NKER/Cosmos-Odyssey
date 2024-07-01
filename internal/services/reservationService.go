package services

import (
	"Cosmos-Odyssey/internal/models"
	"Cosmos-Odyssey/internal/services/external"
	"database/sql"
	"errors"
	"log"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ReservationService struct {
	DB        *sql.DB
	Pricelist *external.PricelistService
}

func NewReservationService(database *sql.DB, pricelistService *external.PricelistService) *ReservationService {
	return &ReservationService{
		DB:        database,
		Pricelist: pricelistService,
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

	_, err = tx.Exec("INSERT INTO reservations (id, pricelist_id, first_name, last_name, routes, total_quoted_price, total_quoted_travel_time, travel_companies, flight_ids) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		reservation.ID, reservation.PricelistID, reservation.Firstname, reservation.Lastname, reservation.Route, reservation.TotalQuotedPrice, reservation.TotalQuotedTravelTime, reservation.TransportationCompanyNames, reservation.RouteIDs)
	if err != nil {
		log.Println("Failed to insert reservation: ", err)
		return reservation, err
	}

	return reservation, nil
}

func (r *ReservationService) createArray(route string, stringType string) []string {
	var newRoute []string
	route = strings.Trim(route, "[]")

	switch stringType {
	case "transport":
		segments := strings.Split(route, ", ")
		for _, seg := range segments {
			if seg != "" {
				newRoute = append(newRoute, seg)
			}
		}

	case "route":
		segments := strings.Split(route, " ")
		for _, seg := range segments {
			seg = strings.TrimSpace(seg)
			if seg != "" {
				newRoute = append(newRoute, seg)
			}
		}
	case "routeID":
		segments := strings.Split(route, ", ")
		for _, seg := range segments {
			if seg != "" {
				newRoute = append(newRoute, seg)
			}
		}
	}

	return newRoute
}

func (r *ReservationService) ValidateReservation(reservation models.Reservation) error {
	parsedTime, err := time.Parse("Jan 02, 2006 15:04 GMT", reservation.ValidUntil)
	if err != nil {
		log.Println("Error parsing reservation.ValidUntil: ", err)
		return err
	}
	timeNow := time.Now().UTC()
	if parsedTime.Before(timeNow) {
		return errors.New("invalid reservation, there's a new pricelist already, refresh and try again")
	}

	routeArray := r.createArray(reservation.Route, "route")
	transportationArray := r.createArray(reservation.TransportationCompanyNames, "transport")
	routeIDArray := r.createArray(reservation.RouteIDs, "routeID")

	r.Pricelist.Mutex.Lock()
	defer r.Pricelist.Mutex.Unlock()

	for i := 0; i < len(routeArray)-1; i++ {
		from := routeArray[i]
		to := routeArray[i+1]
		providersFound := false
		for _, leg := range r.Pricelist.Pricelist.Legs {
			routeInfo := leg.RouteInfo
			if routeInfo.From.Name == from && routeInfo.To.Name == to {
				for _, provider := range leg.Providers {
					if provider.Company.Name == transportationArray[i] {
						providersFound = true
						break
					}
				}
				break
			}
		}

		if !providersFound {
			return errors.New("providers for route segment not found")
		}
	}

	totalQuotedPrice := 0.0
	for i := 0; i < len(routeArray)-1; i++ {
		from := routeArray[i]
		to := routeArray[i+1]

		for _, leg := range r.Pricelist.Pricelist.Legs {
			routeInfo := leg.RouteInfo
			if routeInfo.From.Name == from && routeInfo.To.Name == to {
				for _, provider := range leg.Providers {
					if provider.Id == routeIDArray[i] {
						totalQuotedPrice += provider.Price
						break
					}
				}
				break
			}
		}
	}
	if math.Abs(math.Round(totalQuotedPrice)-reservation.TotalQuotedPrice) > 1 {
		return errors.New("total quoted price does not match")
	}

	return nil
}

func (r *ReservationService) GetReservation(reservation models.Reservation) ([]models.Reservation, error) {
	var reservations []models.Reservation
	query := "SELECT id, first_name, last_name, routes, total_quoted_price, total_quoted_travel_time, travel_companies FROM reservations WHERE first_name = ? AND last_name = ?"

	rows, err := r.DB.Query(query, reservation.Firstname, reservation.Lastname)
	if err != nil {
		log.Println("Error getting reservation: ", err)
		return reservations, err
	}

	defer rows.Close()

	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(&reservation.ID, &reservation.Firstname, &reservation.Lastname, &reservation.Route, &reservation.TotalQuotedPrice, &reservation.TotalQuotedTravelTime, &reservation.TransportationCompanyNames) 
		if err != nil {
			log.Println("Error scanning reservation:", err)
			return reservations, err
		}
		reservations = append(reservations, reservation)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Println("Rows iteration error:", err)
		return reservations, err
	}

	return reservations, nil
}
