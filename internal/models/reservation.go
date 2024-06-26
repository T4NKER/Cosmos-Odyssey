package models

type Reservation struct {
	ID                         string   `form:"id"`
	PricelistID                string   `form:"pricelistID"`
	Firstname                  string   `form:"firstname"`
	Lastname                   string   `form:"lastname"`
	Route                      string   `form:"fullRoute"`
	TotalQuotedPrice           float64  `form:"totalCost"`
	RouteIDs                   string `form:"routeIDs"`
	TotalQuotedTravelTime      string   `form:"totalTime"`
	TransportationCompanyNames string   `form:"transportationCompanyNames"`
	ValidUntil                 string   `form:"validUntil"`
}
