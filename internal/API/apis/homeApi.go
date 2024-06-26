package apis

import (
	"Cosmos-Odyssey/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeAPI struct {
	HomeService *services.HomeService
}

func NewHomeAPI(homeService *services.HomeService) *HomeAPI {
	return &HomeAPI{HomeService: homeService}
}

func (h *HomeAPI) RegisterRoutes(router *gin.Engine) {
	router.GET("/", h.HomeHandler)
}

func (h *HomeAPI) HomeHandler(c *gin.Context) {
    CompaniesAndDestinations, err := h.HomeService.GetDestinationsAndCompanies()
    if err != nil {
        // Handle error if data fetching fails
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "error": "Error getting data, please return to homepage",
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