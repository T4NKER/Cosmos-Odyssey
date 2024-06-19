package models

type Reservation struct {
	ID                         string   `json:"id"`
	PricelistID                string   `json:"pricelistID"`
	Firstname                  string   `json:"firstname"`
	Lastname                   string   `json:"lastname"`
	Route                      []string `json:"route"`
	TotalQuotedPrice           float64  `json:"totalQuotedPrice"`
	TotalQuotedTravelTime      float64  `json:"totalQuotedTravelTime"`
	TransportationCompanyNames string   `json:"transportationCompanyNames"`
}
