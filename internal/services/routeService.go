// The RouteService struct in Go contains methods for retrieving flight plans, filtering routes,
// sorting routes, finding all possible routes, and fetching providers for a given route.
package services

import (
	"Cosmos-Odyssey/internal/models"
	"Cosmos-Odyssey/internal/services/external"
	"database/sql"
	"log"
	"sort"
	"sync"

	"time"
)

type RouteService struct {
	DB            *sql.DB
	Pricelist     models.Pricelist
	pricelistMutex sync.Mutex
}

func NewRouteService(database *sql.DB) *RouteService {
	service := &RouteService{
		DB: database,
	}

	go service.pricelistUpdater()

	return service
}

func (r *RouteService) pricelistUpdater() {
	r.Pricelist = external.ApiScheduler()
	for {
		time.Sleep(time.Minute) // Check every minute; adjust as needed

		r.pricelistMutex.Lock()
		if r.Pricelist.ValidUntil.Before(time.Now()) {
			r.Pricelist = external.ApiScheduler()
			log.Println("Pricelist refreshed.")
		}
		r.pricelistMutex.Unlock()
	}
}

func (r *RouteService) GetQuotes(requestedRoute models.RequestedRoute) ([]models.QuotedRoute, error) {
	startTime := time.Now()
	log.Println(time.Now().UTC())
	log.Println(r.Pricelist.ValidUntil)
	visited := make(map[string]bool)
	currentPath := []string{}
	foundRoutes := [][]string{}

	// Find routes
	foundRoutes = r.findRoutes(requestedRoute.From, requestedRoute.To, external.AllConnections, visited, currentPath)
	log.Println(foundRoutes)
	quotedRoutes, err := r.findAllPossiblePrices(foundRoutes, requestedRoute)
	if err != nil {
		return nil, err
	}

	quotedRoutes, err = r.sorter(quotedRoutes, requestedRoute)
	if err != nil {
		return nil, err
	}

	elapsedTime := time.Since(startTime)
	log.Printf("GetQuotes took %s", elapsedTime)

	return quotedRoutes, nil
}

func (r *RouteService) sorter(quotedRoutes []models.QuotedRoute, requestedRoute models.RequestedRoute) ([]models.QuotedRoute, error) {
	switch requestedRoute.Sort {
	case "price":
		sort.SliceStable(quotedRoutes, func(i, j int) bool {
			return quotedRoutes[i].TotalCost < quotedRoutes[j].TotalCost
		})
	case "time":
		sort.SliceStable(quotedRoutes, func(i, j int) bool {
			return quotedRoutes[i].TotalTime < quotedRoutes[j].TotalTime
		})
	case "distance":
		sort.SliceStable(quotedRoutes, func(i, j int) bool {
			return quotedRoutes[i].TotalDistance < quotedRoutes[j].TotalDistance
		})
	default:
		sort.SliceStable(quotedRoutes, func(i, j int) bool {
			return quotedRoutes[i].TotalCost*1.618*float64(quotedRoutes[i].TotalDistance) < quotedRoutes[j].TotalCost*1.618*float64(quotedRoutes[j].TotalDistance)
		})
	}
	return quotedRoutes, nil
}

func (r *RouteService) findRoutes(origin string, destination string, connections map[string][]string, visited map[string]bool, currentPath []string) [][]string {
    visited[origin] = true
    currentPath = append(currentPath, origin)
    var allPaths [][]string
    if origin == destination {
        newPath := make([]string, len(currentPath))
        copy(newPath, currentPath)
        allPaths = append(allPaths, newPath)
    } else {
        for _, next := range connections[origin] {
            if !visited[next] {
                paths := r.findRoutes(next, destination, connections, visited, currentPath)
                allPaths = append(allPaths, paths...)
            }
        }
    }
    visited[origin] = false
    currentPath = currentPath[:len(currentPath)-1]
    return allPaths
}

func (r *RouteService) findAllPossiblePrices(routes [][]string, requestedRoute models.RequestedRoute) ([]models.QuotedRoute, error) {
	var allQuotedRoutes []models.QuotedRoute
	for _, route := range routes {
		quotedRoutes := r.generateRoutesForProviders(route, requestedRoute, models.QuotedRoute{FullRoute: route}, 0)
		allQuotedRoutes = append(allQuotedRoutes, quotedRoutes...)
	}
	return allQuotedRoutes, nil
}

func (r *RouteService) generateRoutesForProviders(route []string, requestedRoute models.RequestedRoute, currentQuotedRoute models.QuotedRoute, segmentIndex int) []models.QuotedRoute {
	if segmentIndex >= len(route)-1 {
		return []models.QuotedRoute{currentQuotedRoute}
	}

	from := route[segmentIndex]
	to := route[segmentIndex+1]

	providers, err := r.fetchProviders(from, to, requestedRoute)
	if err != nil {
		log.Println("Error fetching providers:", err)
		return nil
	}

	var allQuotedRoutes []models.QuotedRoute

	for _, provider := range providers {
		if segmentIndex > 0 {
			previousProvider := currentQuotedRoute.Sections[segmentIndex-1]
			if provider.FlightStart.Before(previousProvider.FlightEnd) {
				continue
			}
		}

		newQuotedRoute := models.QuotedRoute{
			FullRoute:     currentQuotedRoute.FullRoute,
			Sections:      append([]models.RouteSection{}, currentQuotedRoute.Sections...),
			TotalCost:     currentQuotedRoute.TotalCost + provider.Price,
			TotalDistance: currentQuotedRoute.TotalDistance + provider.Distance,
		}

		newQuotedRoute.Sections = append(newQuotedRoute.Sections, provider)

		var journeyStart time.Time
		if segmentIndex == 0 {
			journeyStart = provider.FlightStart
		} else {
			journeyStart = currentQuotedRoute.Sections[0].FlightStart
		}
		journeyEnd := provider.FlightEnd

		newQuotedRoute.TotalTime = journeyEnd.Sub(journeyStart)

		furtherRoutes := r.generateRoutesForProviders(route, requestedRoute, newQuotedRoute, segmentIndex+1)
		allQuotedRoutes = append(allQuotedRoutes, furtherRoutes...)
	}

	return allQuotedRoutes
}

func (r *RouteService) fetchProviders(from, to string, requestedRoute models.RequestedRoute) ([]models.RouteSection, error) {
    var providers []models.RouteSection
	currentTime := time.Now().UTC()
    for _, leg := range r.Pricelist.Legs {
        routeInfo := leg.RouteInfo
        if routeInfo.From.Name == from && routeInfo.To.Name == to {
            for _, provider := range leg.Providers {
                if provider.FlightStart.After(currentTime) && requestedRoute.Provider == "" || provider.FlightStart.After(currentTime) && provider.Company.Name == requestedRoute.Provider {
                    routeSection := models.RouteSection{
                        ID:          provider.Id,
                        Provider:    provider.Company.Name,
                        Price:       provider.Price,
                        Distance:    routeInfo.Distance,
                        FlightStart: provider.FlightStart,
                        FlightEnd:   provider.FlightEnd,
                        From:        from,
                        To:          to,
                        Time:        provider.FlightEnd.Sub(provider.FlightStart).Hours(),
                    }
                    providers = append(providers, routeSection)
                }
            }
        }
    }

    return providers, nil
}
