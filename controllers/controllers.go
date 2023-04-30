package controllers

import (
	"bd_aerolinea/database"
	"bd_aerolinea/models"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GetVuelo is a handler function that retrieves a flight based on its origin, destination and date of departure.
func GetVuelos(c *gin.Context) {

	origenVuelo := c.Query("origen")
	destinoVuelo := c.Query("destino")
	fechaVuelo := c.Query("fecha")

	// Call the GetVuelo function from the models package to retrieve the flight
	vuelos, err := database.GetVuelos(origenVuelo, destinoVuelo, fechaVuelo)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)

	} else {
		c.JSON(http.StatusOK, vuelos)
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
		database.CreateVuelo(vuelo)
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
	response, err := database.UpdateVuelo(numeroVuelo, origenVuelo, destinoVuelo, fechaVuelo, update.StockDePasajeros)

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
	err := database.DeleteVuelo(numeroVuelo, origenVuelo, destinoVuelo, fechaVuelo)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.Status(http.StatusNoContent)
	}
}

func CreateReservation(c *gin.Context) {

	var reserva models.Reservation

	// Bind the request body to the reserva variable
	if err := c.BindJSON(&reserva); err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
	} else {

		fmt.Println("Creating Reservation")

		// Call the CreateReservation function from the models package to create the new reservation
		response, err := database.CreateReservation(reserva)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		c.IndentedJSON(http.StatusOK, response)
	}
}

func GetReservations(c *gin.Context) {

	// Retrieve the reservation information (pnr, apellido) from the query parameters
	pnr := c.Query("pnr")
	apellido := c.Query("apellido")

	// Call the GetReservation func to get the said reservation using the parameters
	reserva, err := database.GetReservation(pnr, apellido)

	if err != nil {
		fmt.Println(err)

		c.IndentedJSON(http.StatusNotFound, strings.TrimSuffix(fmt.Sprintln(err), "\n"))

	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"vuelos": reserva.Vuelos, "pasajeros": reserva.Passengers})
	}
}

func UpdateReservation(c *gin.Context) {

	// Retrieve the reservation information (pnr, apellido) from the query parameters
	pnr := c.Query("pnr")
	apellido := c.Query("apellido")

	var reserva models.Reservation

	err := c.BindJSON(&reserva)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Reseting the pnr and surname
	reserva.PNR = pnr
	reserva.Apellido = apellido

	// Call the UpdateReservation function from the models package to remplace the reservation with a other one.
	response, err := database.UpdateReservation(pnr, apellido, reserva)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, response)

}

func DeleteReservation(c *gin.Context) {
	// Retrieve the reservation information (pnr, apellido) from the query parameters
	pnr := c.Query("pnr")
	apellido := c.Query("apellido")

	// Call the DeleteReservation function from the models package to delete the reservation with the given parameters
	err := database.DeleteReservation(pnr, apellido)

	if err != nil {
		fmt.Println(err)

		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.Status(http.StatusNoContent)
	}
}

func Max(data map[string]int) string {
	var max int = 0
	var maxIndex string
	for i, v := range data {
		if v > max {
			maxIndex = i
			max = v
		}
	}
	return maxIndex
}

func Min(data map[string]int) string {
	var min int = 2000000000
	var minIndex string
	for i, v := range data {
		if v < min {
			minIndex = i
			min = v
		}
	}
	return minIndex
}

func GetStatistics(c *gin.Context) {
	var (
		balancesAncillaries = make(map[string]int)
		stats               models.Statistics
		RutaGanancia        = make(map[string]int)
	)
	reservas, err := database.GetAllReservations()

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	for i := 0; i < len(reservas); i++ {
		RutaStringIda := reservas[i].Vuelos[0].Origen + " - " + reservas[i].Vuelos[0].Destino

		for j := 0; j < len(reservas[i].Passengers); j++ {
			RutaGanancia[RutaStringIda] += reservas[i].Passengers[j].Balances.VueloIda

			for k := 0; k < len(reservas[i].Passengers[j].Ancillaries.Ida); k++ {
				balancesAncillaries[reservas[i].Passengers[j].Ancillaries.Ida[k].SSR] += reservas[i].Passengers[j].Balances.AncillariesIda
			}

			if len(reservas[i].Vuelos) == 2 {
				RutaStringVuelta := reservas[i].Vuelos[1].Origen + " - " + reservas[i].Vuelos[1].Destino
				RutaGanancia[RutaStringVuelta] += reservas[i].Passengers[j].Balances.VueloVuelta

				for k := 0; k < len(reservas[i].Passengers[j].Ancillaries.Vuelta); k++ {
					balancesAncillaries[reservas[i].Passengers[j].Ancillaries.Vuelta[k].SSR] += reservas[i].Passengers[j].Balances.AncillariesVuelta
				}
			}

		}

	}
	stats.RutaMayorGanancia = Max(RutaGanancia)
	stats.RutaMenorGanancia = Min(RutaGanancia)
	for i, v := range models.AncillaryName {

		var RankingAncillary models.AncillariesStatistics
		RankingAncillary.Nombre = v
		RankingAncillary.SSR = i
		RankingAncillary.Ganancia = balancesAncillaries[i]
		stats.RankingAncillaries = append(stats.RankingAncillaries, RankingAncillary)
	}

	layout := "02/01/2006"

	var pasajerosMes map[string]int = map[string]int{"January": 0, "February": 0, "March": 0, "April": 0, "May": 0, "June": 0, "July": 0, "August": 0, "September": 0, "October": 0, "November": 0, "December": 0}

	for i := 0; i < len(reservas); i++ {
		dateIda := reservas[i].Vuelos[0].Fecha
		mesIda, _ := time.Parse(layout, dateIda)
		pasajerosMes[fmt.Sprint(mesIda.Month())] += len(reservas[i].Passengers)
		if len(reservas[i].Vuelos) == 2 {
			dateVuelta := reservas[i].Vuelos[1].Fecha
			mesVuelta, _ := time.Parse(layout, dateVuelta)
			pasajerosMes[fmt.Sprint(mesVuelta.Month())] += len(reservas[i].Passengers)
		}
	}

	stats.PromedioPasajeros.Jan = pasajerosMes["January"]
	stats.PromedioPasajeros.Feb = pasajerosMes["February"]
	stats.PromedioPasajeros.Mar = pasajerosMes["March"]
	stats.PromedioPasajeros.Apr = pasajerosMes["April"]
	stats.PromedioPasajeros.May = pasajerosMes["May"]
	stats.PromedioPasajeros.Jun = pasajerosMes["June"]
	stats.PromedioPasajeros.Jul = pasajerosMes["July"]
	stats.PromedioPasajeros.Aug = pasajerosMes["August"]
	stats.PromedioPasajeros.Sep = pasajerosMes["September"]
	stats.PromedioPasajeros.Oct = pasajerosMes["October"]
	stats.PromedioPasajeros.Nov = pasajerosMes["November"]
	stats.PromedioPasajeros.Dec = pasajerosMes["December"]

	c.IndentedJSON(http.StatusOK, stats)
}

func GenerateNewPNR(c *gin.Context) {
	reservations, err := database.GetAllReservations()
	if err != nil {
		log.Fatal("Error in getting the reservations")
	}
	for true {
		//The max number of pnr combinations is 36^6
		// In this function we are generating a pnr then looking if we already have a
		// reservation with this pnr if yes we generate a new pnr, if not we continue
		// with the generated pnr.
		n := rand.Intn(36 * 36 * 36 * 36 * 36 * 36)
		NewPnr := convertToBase(n, 36)
		if elementInReservations(NewPnr, reservations) == 0 {
			var PNRC models.PNRCapsule
			PNRC.PNR = NewPnr
			c.JSON(http.StatusOK, PNRC)
			break
		}
	}
}

func elementInReservations(NewPNR string, reservations []models.Reservation) int {
	//This function looks in a array of reservations if the inputed pnr is already present
	// in an other reservation.
	if reservations == nil {
		return 0
	}
	for i := 0; i < len(reservations); i++ {
		if reservations[i].PNR == NewPNR {
			return 1
		}
	}
	return 0
}

func convertToBase(n, base int) string {
	//This function is used to converrt an int in base 10 to any other base
	// A pnr is a combination of the letters and numbers
	// There are 36 in total so we generate a number between 0 and 36^6
	// Then we convert it in base 36 to get a randmly generated pnr

	if n == 0 {
		return "0"
	}
	digits := []string{}
	for n > 0 {
		modulo := n % base
		if modulo < 10 {
			digits = append(digits, fmt.Sprintf("%d", modulo))
		} else {
			digits = append(digits, string('A'-10+modulo))
		}
		n /= base
	}
	reverse(digits)
	return strings.Join(digits, "")
}

func reverse(a []string) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}
