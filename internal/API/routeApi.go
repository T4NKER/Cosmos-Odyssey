package apis

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


func GetRoutes(c *gin.Context) {
    origin := c.Query("origin")
    destination := c.Query("destination")
    var routes []models.Route

    

    c.JSON(http.StatusOK, routes)
}
