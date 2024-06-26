package models

import "time"

// "time"

type RequestedRoute struct {
	To       string `form:"to" binding:"required"`
	From     string `form:"from" binding:"required"`
	Provider string `form:"provider"`
	Sort     string `form:"sort"`
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
	MaxCost  int    `form:"max_cost"`
}

type QuotedRoute struct {
	PricelistID   string
	FullRoute     []string
	Sections      []RouteSection
	TotalCost     int64
	TotalTime     string
	TotalDistance int
	ValidUntil    time.Time
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
