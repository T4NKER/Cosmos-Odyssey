package external

import (
	"Cosmos-Odyssey/internal/models"
	database "Cosmos-Odyssey/pkg/database"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	//"github.com/go-co-op/gocron/v2"
)

var AllConnections = make(map[string][]string)

func ApiScheduler() (models.Pricelist) {
	pricelistData := fetchAndStoreData()
	return pricelistData
}

func fetchAndStoreData() (models.Pricelist) {
	response, err := http.Get("https://cosmos-odyssey.azurewebsites.net/api/v1.0/TravelPrices")
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}
	defer response.Body.Close()

	var pricelistData models.Pricelist
	if err := json.NewDecoder(response.Body).Decode(&pricelistData); err != nil {
		log.Fatalf("Failed to decode data: %v", err)
	}

	// Store pricelist data in the database
	storePricelist(database.Db, pricelistData)

	// Update in-memory connections
	updateConnections(pricelistData)

	return pricelistData
}

func storePricelist(db *sql.DB, pricelist models.Pricelist) {
	tx, err := db.Begin()
	if err != nil {
		log.Println("Failed to begin transaction: ", err)
	}

	_, err = tx.Exec("INSERT INTO pricelists (id, valid_until) VALUES ($1, $2)", pricelist.Id, pricelist.ValidUntil)
	if err != nil {
		tx.Rollback()
		log.Println("Failed to insert pricelist: ", err)
	}

	for _, route := range pricelist.Legs {
		_, err := tx.Exec("INSERT INTO routes (id, pricelist_id, from_planet, to_planet, distance) VALUES ($1, $2, $3, $4, $5)",
			route.Id, pricelist.Id, route.RouteInfo.From.Name, route.RouteInfo.To.Name, route.RouteInfo.Distance)
		if err != nil {
			tx.Rollback()
			log.Println("Failed to insert route: ", err)
		}

		for _, provider := range route.Providers {
			_, err := tx.Exec("INSERT INTO providers (id, route_id, company_name, price, flight_start, flight_end) VALUES ($1, $2, $3, $4, $5, $6)",
				provider.Id, route.Id, provider.Company.Name, provider.Price, provider.FlightStart, provider.FlightEnd)
			if err != nil {
				tx.Rollback()
				log.Println("Failed to insert provider: ", err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Failed to commit transaction: ", err)
		log.Println("This is probably because of trying to refresh the flight data too early and thus the unique id being the same.")
	}

	_, err = db.Exec("DELETE FROM pricelists WHERE id NOT IN (SELECT id FROM pricelists ORDER BY created_at DESC LIMIT 15)")
	if err != nil {
		log.Println("Failed to delete old pricelists: ", err)
	}
}

func updateConnections(pricelist models.Pricelist) {
	for _, route := range pricelist.Legs {
		fromPlanet := route.RouteInfo.From.Name
		toPlanet := route.RouteInfo.To.Name

		if _, ok := AllConnections[fromPlanet]; !ok {
			AllConnections[fromPlanet] = []string{}
		}

		AllConnections[fromPlanet] = append(AllConnections[fromPlanet], toPlanet)
	}
}
