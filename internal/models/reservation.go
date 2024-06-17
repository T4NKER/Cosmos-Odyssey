package models

type Reservation struct {
	ID                        string   `json:"id"`
	Firstname                 string   `json:"firstname"`
	Lastname                  string   `json:"lastname"`
	Route                     []string `json:"route"`
	TotalQuotedPrice          float64  `json:"totalQuotedPrice"`
	TotalQuotedTravelTime     float64  `json:"totalQuotedTravelTime"`
	TransportationCompanyName []string `json:"transportationCompanyName"`
}
