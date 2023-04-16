package main

import (
	"bd_aerolinea/controllers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load the environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var port = os.Getenv("PORT")
	var server = os.Getenv("SERVER")
	router := gin.Default()

	router.GET("/api/vuelo", controllers.GetVuelo)
	router.PUT("/api/vuelo", controllers.UpdateVuelo)
	router.DELETE("/api/vuelo", controllers.DeleteVuelo)
	router.POST("/api/vuelo", controllers.CreateVuelo)

	router.Run(server + ":" + port)
}
