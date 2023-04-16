package controllers

import (
	"bd_aerolinea/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*func GetVuelos(c *gin.Context) {

	//vuelo := models.GetVuelos()

	if vuelo == nil || len(vuelo) == 0 {

		c.AbortWithStatus(http.StatusNotFound)

	} else {

		c.IndentedJSON(http.StatusOK, vuelo)

	}
}*/

// GetVuelo is a handler function that retrieves a flight based on its origin, destination and date of departure.
func GetVuelos(c *gin.Context) {

	origenVuelo := c.Query("origen")
	destinoVuelo := c.Query("destino")
	fechaVuelo := c.Query("fecha")

	// Call the GetVuelo function from the models package to retrieve the flight
	vuelos, err := models.GetVuelos(origenVuelo, destinoVuelo, fechaVuelo)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)

	} else {
		c.IndentedJSON(http.StatusOK, gin.H{
			"vuelos": vuelos,
		})
	}

}

// CreateVuelo is a handler function that creates a new flight.
func CreateVuelo(c *gin.Context) {

	var vuelo models.Flight

	// Bind the request body to the vuelo variable
	if err := c.BindJSON(&vuelo); err != nil {

		c.AbortWithStatus(http.StatusBadRequest)
	} else {

		// Call the CreateVuelo function from the models package to create the new flight
		models.CreateVuelo(vuelo)
	}
}

// UpdateVuelo is a handler function that updates an existing flight.
func UpdateVuelo(c *gin.Context) {

	//var vuelo models.Flight

	// Retrieve the flight information (origin, destination, date, and flight number) from the query parameters
	numeroVuelo := c.Query("numero_vuelo")
	origenVuelo := c.Query("origen")
	destinoVuelo := c.Query("destino")
	fechaVuelo := c.Query("fecha")

	var update struct {
		StockDePasajeros int `json:"stock_de_pasajeros"`
	}

	err := c.ShouldBindJSON(&update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the UpdateStockPasajeros function from the models package to update the flight's passenger stock
	response, err := models.UpdateVuelo(numeroVuelo, origenVuelo, destinoVuelo, fechaVuelo, update.StockDePasajeros)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteVuelo is a handler function that deletes an existing flight.
func DeleteVuelo(c *gin.Context) {

	// Retrieve the flight information (origin, destination, date, and flight number) from the query parameters
	numeroVuelo := c.Query("numero_vuelo")
	origenVuelo := c.Query("origen")
	destinoVuelo := c.Query("destino")
	fechaVuelo := c.Query("fecha")

	// Call the DeleteVuelo function from the models package to delete the flight with the given parameters
	err := models.DeleteVuelo(numeroVuelo, origenVuelo, destinoVuelo, fechaVuelo)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.Status(http.StatusNoContent)
	}
}
