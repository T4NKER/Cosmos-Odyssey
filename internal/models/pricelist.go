package models

import (
	"time"
)

// This should be refetched every 15 minutes, since the data refreshes every 15 minutes.
type Pricelist struct {
	Id         string `json:"id"`
	ValidUntil string `json:"validUntil"`
	Legs       []Legs `json:"legs"`
}

type Legs struct {
	Id        string      `json:"id"`
	RouteInfo RouteInfo   `json:"routeInfo"`
	Providers []Providers `json:"providers"`
}

type RouteInfo struct {
	Id   string `json:"id"`
	From struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"from"`
	To struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"to"`
	Distance int `json:"distance"`
}

type Providers struct {
	Id      string `json:"id"`
	Company struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"company"`
	Price       float64   `json:"price"`
	FlightStart time.Time `json:"flightStart"`
	FlightEnd   time.Time `json:"flightEnd"`
}
