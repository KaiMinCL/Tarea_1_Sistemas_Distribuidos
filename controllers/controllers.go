package controllers

import (
	"bd_aerolinea/models"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"time"
)

var(
	AncillaryPrice = map[string]int{"BGH":10000, "BGR":30000, "STDF":5000, "PAXS":2000, "PTCR":40000, "AVIH":40000, "SPML":35000, "LNGE":15000, "WIFI":20000}
)

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

func GenerateNewPNR() string {
	reservations, err := models.GetAllReservations()

	if err != nil {
		log.Fatal("Error in getting the reservations")
	}

	if reservations == nil {
		return convertToBase(rand.Intn(36*36*36*36*36*36), 36)
	}

	for true {
		//The max number of pnr combinations is 36^6
		// In this function we are generating a pnr then looking if we already have a
		// reservation with this pnr if yes we generate a new pnr, if not we continue
		// with the generated pnr.
		n := rand.Intn(36 * 36 * 36 * 36 * 36 * 36)
		NewPnr := convertToBase(n, 36)
		if elementInReservations(NewPnr, reservations) == 0 {
			return NewPnr
		}
	}

	return "This return wil never be executed"
}

func elementInReservations(NewPNR string, reservations []models.Reservation) int {
	//This function looks in a array of reservations if the inputed pnr is already present
	// in an other reservation.
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

func SumAncillaries(PA models.PassengerAncillaryList) (int, int){
	//This function calculates the price for the ancillaries associated with this passenger.
	countIda := 0
	countVuelta := 0
	for j := 0; j < len(PA.Ida); j++{
		countIda += AncillaryPrice[PA.Ida[j].SSR]*PA.Ida[j].Cantidad
	}
	for j := 0; j < len(PA.Vuelta); j++{
		countVuelta += AncillaryPrice[PA.Vuelta[j].SSR]*PA.Vuelta[j].Cantidad
		fmt.Println(countVuelta, PA.Vuelta[j].SSR, PA.Vuelta[j].Cantidad)
	}

	return countIda, countVuelta
}

func SumVuelos (flights []models.ReservationFlight) (int, int){
	//This function returns the sum of the price of all the tickets associated with this passenger.
	var tiempoIda int = 0
	var tiempoVuelta int = 0

	horaSalidaIda, _ := time.Parse("15:04", flights[0].HoraSalida)
	horaLlegadaIda, _ := time.Parse("15:04", flights[0].HoraLlegada)

	if horaLlegadaIda.Before(horaSalidaIda) {
		horaLlegadaIda = horaLlegadaIda.Add(24 * time.Hour)
	}

	tiempoIda = int(horaLlegadaIda.Sub(horaSalidaIda).Minutes())

	if len(flights) == 2{
		horaSalidaVuelta, _ :=time.Parse("15:04", flights[1].HoraSalida)
		horaLlegadaVuelta, _ :=time.Parse("15:04", flights[1].HoraLlegada)

		if horaLlegadaVuelta.Before(horaSalidaVuelta) {
			horaLlegadaVuelta = horaLlegadaVuelta.Add(24 * time.Hour)
		}

		tiempoVuelta = int(horaLlegadaVuelta.Sub(horaSalidaVuelta).Minutes())
	} else if len(flights) > 2{
		fmt.Print("There shouldn't be more than tow flights per reservation")
	}

	return tiempoIda*590, tiempoVuelta*590
}

func CreateReservation(c *gin.Context) {

	var reserva models.Reservation

	// Bind the request body to the reserva variable
	if err := c.BindJSON(&reserva); err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
	} else {

		fmt.Println("Calculating Balances")
		for i := 0; i<len(reserva.Passengers); i++{
			reserva.Passengers[i].Balances.AncillariesIda, reserva.Passengers[i].Balances.AncillariesVuelta = SumAncillaries(reserva.Passengers[i].Ancillaries)
			reserva.Passengers[i].Balances.VueloIda, reserva.Passengers[i].Balances.VueloVuelta = SumVuelos(reserva.Vuelos)
		}


		fmt.Println("Creating Reservation")

		// Call the CreateReservation function from the models package to create the new reservation
		response, err := models.CreateReservation(reserva)
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
	reservas, err := models.GetReservation(pnr, apellido)

	if err != nil {
		fmt.Println(err)

		c.IndentedJSON(http.StatusNotFound, strings.TrimSuffix(fmt.Sprintln(err), "\n"))

	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"reservas": reservas})
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

	//Calculating the new balances
	for i := 0; i<len(reserva.Passengers); i++{
			reserva.Passengers[i].Balances.AncillariesIda, reserva.Passengers[i].Balances.AncillariesVuelta = SumAncillaries(reserva.Passengers[i].Ancillaries)
			reserva.Passengers[i].Balances.VueloIda, reserva.Passengers[i].Balances.VueloVuelta = SumVuelos(reserva.Vuelos)
		}

	// Call the UpdateReservation function from the models package to remplace the reservation with a other one.
	response, err := models.UpdateReservation(pnr, apellido, reserva)

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
	err := models.DeleteReservation(pnr, apellido)

	if err != nil {
		fmt.Println(err)

		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.Status(http.StatusNoContent)
	}
}
