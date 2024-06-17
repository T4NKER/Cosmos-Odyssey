package apis

import (
	"Cosmos-Odyssey/internal/models"
	"Cosmos-Odyssey/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouteAPI struct {
    RouteService *services.RouteService
}

func NewRouteAPI(routeService *services.RouteService) *RouteAPI {
	return &RouteAPI{RouteService: routeService}
}

func (r *RouteAPI) RegisterRoutes(router *gin.Engine) {
    router.POST("/routes", r.GetQuotes)
}

func (r *RouteAPI) GetQuotes(c *gin.Context) {
    var route models.RequestedRoute
	if err := c.BindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

    quotedRoutes, err  := r.RouteService.GetQuotes(route) 
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }

    c.JSON(http.StatusOK, quotedRoutes)
}
