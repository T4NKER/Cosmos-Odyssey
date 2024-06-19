package services

import (
	"Cosmos-Odyssey/internal/models"
	"database/sql"
	"fmt"
)

type HomeService struct {
	DB            *sql.DB
}

func NewHomeService(DB *sql.DB) *HomeService {
	return &HomeService{
		DB: DB,
	}
}

func (h *HomeService) GetDestinationsAndCompanies() (models.CompaniesAndDestinations, error) {
	var result models.CompaniesAndDestinations
	companyRows, err := h.DB.Query("SELECT DISTINCT company_name FROM providers")
	if err != nil {
		return result, fmt.Errorf("failed to query companies: %v", err)
	}
	defer companyRows.Close()

	destinationRows, err := h.DB.Query("SELECT DISTINCT from_planet FROM routes")
	if err != nil {
		return result, fmt.Errorf("failed to query destinations: %v", err)
	}
	defer destinationRows.Close()

	for companyRows.Next() {
		var companyName string
		if err := companyRows.Scan(&companyName); err != nil {
			return result, fmt.Errorf("failed to scan company name: %v", err)
		}
		result.Companies = append(result.Companies, companyName)
	}

	for destinationRows.Next() {
		var destinationName string
		if err := destinationRows.Scan(&destinationName); err != nil {
			return result, fmt.Errorf("failed to scan destination name: %v", err)
		}
		result.Destinations = append(result.Destinations, destinationName)
	}

	if err := companyRows.Err(); err != nil {
		return result, fmt.Errorf("error encountered during company rows iteration: %v", err)
	}

	if err := destinationRows.Err(); err != nil {
		return result, fmt.Errorf("error encountered during destination rows iteration: %v", err)
	}

	return result, nil
}