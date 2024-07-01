package external

import (
	"Cosmos-Odyssey/internal/models"
	database "Cosmos-Odyssey/pkg/database"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
	//"github.com/go-co-op/gocron/v2"
)

var AllConnections = make(map[string][]string)
var pricelistService *PricelistService

type PricelistService struct {
	Pricelist models.Pricelist
	Mutex     sync.Mutex
	DB        *sql.DB
}

func NewExternalPricelistService(DB *sql.DB) *PricelistService {
	pricelistService = &PricelistService{
		Pricelist: models.Pricelist{},
		DB:        DB,
		Mutex:     sync.Mutex{},
	}
	pricelistService.refreshPricelist()
	go pricelistService.pricelistUpdater()

	return pricelistService
}

func (p *PricelistService) pricelistUpdater() {
	for {
		p.Mutex.Lock()
		duration := time.Until(p.Pricelist.ValidUntil)
		log.Println("Setting new duration to ", duration)
		p.Mutex.Unlock()

		timer := time.NewTimer(duration)
		log.Println("Timer set for duration: ", duration)
		<-timer.C
		log.Println("Timer reached end, getting API response")

		log.Println("Refreshing pricelist...")
		p.refreshPricelist()
		log.Println("Pricelist refreshed.")
	}
}

func (p *PricelistService) refreshPricelist() {
	log.Println("Fetching data...")
	pricelistData := fetchAndStoreData()
	log.Println("Data fetched, updating pricelist...")
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	log.Println("Making new pricelist...")
	p.Pricelist = pricelistData
	log.Println("New pricelist made...")
}

func fetchAndStoreData() models.Pricelist {
	log.Println("Starting fetchAndStoreData")
	response, err := http.Get("https://cosmos-odyssey.azurewebsites.net/api/v1.0/TravelPrices")
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}
	defer response.Body.Close()

	var pricelistData models.Pricelist
	log.Println("Decoding response...")
	if err := json.NewDecoder(response.Body).Decode(&pricelistData); err != nil {
		log.Fatalf("Failed to decode data: %v", err)
	}

	log.Println("Storing pricelist data in the database...")
	storePricelist(database.Db, pricelistData)

	updateConnections(pricelistData)
	log.Println("Pricelist data fetched and stored successfully")
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

	_, err = db.Exec(`DELETE FROM pricelists
WHERE id NOT IN (
    SELECT id
    FROM (
        SELECT id
        FROM pricelists
        ORDER BY created_at DESC
        LIMIT 15
    )
);
`)
	if err != nil {
		log.Println("Failed to delete old pricelists: ", err)
	}
}

func updateConnections(pricelist models.Pricelist) {
	AllConnections = make(map[string][]string)

	for _, route := range pricelist.Legs {
		fromPlanet := route.RouteInfo.From.Name
		toPlanet := route.RouteInfo.To.Name

		if _, ok := AllConnections[fromPlanet]; !ok {
			AllConnections[fromPlanet] = []string{}
		}

		AllConnections[fromPlanet] = append(AllConnections[fromPlanet], toPlanet)
	}
}
