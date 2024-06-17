package models

import "time"

// "time"

type RequestedRoute struct {
	To       string `json:"to"`
	From     string `json:"from"`
	Provider string `json:"provider"`
	Sort     string `json:"sort"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	MaxCost  int    `json:"max_cost"`
}

type QuotedRoute struct {
	FullRoute     []string
	Sections      []RouteSection
	TotalCost     float64
	TotalTime     time.Duration
	TotalDistance int
}

type RouteSection struct {
	ID          string
	From        string
	To          string
	Distance    int
	Provider    string
	FlightStart time.Time
	FlightEnd   time.Time
	Price       float64
	Time        float64
}
