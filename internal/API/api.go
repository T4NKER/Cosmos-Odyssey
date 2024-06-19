package API

import (
	"Cosmos-Odyssey/internal/API/apis"
	"Cosmos-Odyssey/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MasterAPI struct {
    RoutesAPI      *apis.RouteAPI
    ReservationAPI *apis.ReservationAPI
    HomeService *services.HomeService
}

func NewMasterAPI(routeAPI *apis.RouteAPI, reservationAPI *apis.ReservationAPI, homeService *services.HomeService) *MasterAPI {
    return &MasterAPI{
        RoutesAPI:      routeAPI,
        ReservationAPI: reservationAPI,
        HomeService: homeService,
    }
}

func (m *MasterAPI) RegisterRoutes(router *gin.Engine) {
    m.RoutesAPI.RegisterRoutes(router)
    m.ReservationAPI.RegisterRoutes(router)
    router.GET("/", m.homeHandler)
}

func (m *MasterAPI) homeHandler(c *gin.Context) {
    CompaniesAndDestinations, err := m.HomeService.GetDestinationsAndCompanies()
    if err != nil {
        // Handle error if data fetching fails
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "message": "Error getting data",
        })
        return
    }

    // Prepare the data to be sent to the HTML template
    data := gin.H{
        "message": "Welcome to Cosmos Odyssey!",
        "companiesAndDestinations": CompaniesAndDestinations,
    }

    // Render the HTML template with the data
    c.HTML(http.StatusOK, "index.html", data)
}