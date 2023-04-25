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

	// Create a new router using the gin framework
	router := gin.Default()

	// Defines the routes for the vuelo API
	router.GET("/api/vuelo", controllers.GetVuelos)
	router.PUT("/api/vuelo", controllers.UpdateVuelo)
	router.DELETE("/api/vuelo", controllers.DeleteVuelo)
	router.POST("/api/vuelo", controllers.CreateVuelo)

	router.GET("/api/reserva", controllers.GetReservations)
	router.PUT("/api/reserva", controllers.UpdateReservation)
	router.DELETE("/api/reserva", controllers.DeleteReservation)
	router.POST("/api/reserva", controllers.CreateReservation)

	router.GET("/api/estadisticas", controllers.GetStatistics)

	router.GET("/api/generatepnr", controllers.GenerateNewPNR)
	router.Run(server + ":" + port)
}
