package server

import (
	apis "Cosomos-Odyssey/internal/API"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
    AllowMethods:     []string{"PUT", "GET", "POST", "DELETE", "OPTIONS"},
	AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
    AllowCredentials: true,
  }))
	apis.Api(router)
	router.Run()

}

