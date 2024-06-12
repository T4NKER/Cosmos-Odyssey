package apis

import (
	"Cosomos-Odyssey/internal/models"
	"Cosomos-Odyssey/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouteAPI struct {
    RouteService *services.RouteService
}

func NewRouteAPI(routeService *RouteService) *RouteAPI {
	return &RouteAPI{RouteService: routeService}
}

func (r *RouteAPI) RegisterRoutes(router *gin.Engine) {
    router.GET("/routes", r.GetRoutes)
}

func (r *RouteAPI) GetRoutes(c *gin.Context) {
    origin := c.Query("origin")
    destination := c.Query("destination")
    var routes []models.Route

    

    c.JSON(http.StatusOK, routes)
}
