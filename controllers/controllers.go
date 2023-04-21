package controllers

import (
	"bd_aerolinea/models"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"strconv"
	"github.com/gin-gonic/gin"
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

func SumAncillaries(PA []models.PassengerAncillaryList) (int, int){
	//This function calculates the price for the ancillaries associated with this passenger.
	countIda := 0
	countVuelta := 0
	for i := 0; i < len(PA); i++{
		for j := 0; j < len(PA[i].Ida); j++{
			countIda += AncillaryPrice[PA[i].Ida[j].SSR]*PA[i].Ida[j].Cantidad
		}
		for j := 0; j < len(PA[i].Vuelta); j++{
			countVuelta += AncillaryPrice[PA[i].Vuelta[j].SSR]*PA[i].Vuelta[j].Cantidad
			fmt.Println(countVuelta, PA[i].Vuelta[j].SSR, PA[i].Vuelta[j].Cantidad)
		}
	}

	return countIda, countVuelta
}

func SumVuelos (c *gin.Context, flights []models.ReservationFlight) (int, int){
	//This function returns the sum of the price of all the tickets associated with this passenger.
	var tiempoIda int
	var tiempoVuelta int

	//We check if we have an even number of flights
	if len(flights) % 2 != 0{
		c.IndentedJSON(http.StatusBadRequest, "The number of flights has to be even")
	}

	//Must be a better way to right this but I don't know how
	//converting the string into int and then calculating the time of the flight to then return the price
	for i:=0; i<len(flights); i+=2{
		horasSalidaIda, minutosSalidaIda, _ := strings.Cut(flights[i].HoraSalida, ":")
		horasLlegadaIda, minutosLlegadaIda, _ := strings.Cut(flights[i].HoraLlegada, ":")
		horasSalidaVuelta, minutosSalidaVuelta, _ := strings.Cut(flights[i+1].HoraSalida, ":")
		horasLlegadaVuelta, minutosLlegadaVuelta, _ := strings.Cut(flights[i+1].HoraLlegada, ":")

		horasSalidaIdaC, _ := strconv.Atoi(horasSalidaIda)
		horasLlegadaIdaC, _ := strconv.Atoi(horasLlegadaIda)
		minutosSalidaIdaC, _ := strconv.Atoi(minutosSalidaIda)
		minutosLlegadaIdaC, _ := strconv.Atoi(minutosLlegadaIda)
		horasSalidaVueltaC, _ := strconv.Atoi(horasSalidaVuelta)
		horasLlegadaVueltaC, _ := strconv.Atoi(horasLlegadaVuelta)
		minutosSalidaVueltaC, _ := strconv.Atoi(minutosSalidaVuelta)
		minutosLlegadaVueltaC, _ := strconv.Atoi(minutosLlegadaVuelta)

		tiempoIda = (1440 * CompareInt(horasSalidaIdaC, horasLlegadaIdaC)) + horasLlegadaIdaC*60+minutosLlegadaIdaC-(horasSalidaIdaC*60+minutosSalidaIdaC)
		tiempoVuelta = (1440 * CompareInt(horasSalidaVueltaC, horasLlegadaVueltaC)) + horasLlegadaVueltaC*60+minutosLlegadaVueltaC-(horasSalidaVueltaC*60+minutosSalidaVueltaC)

	}

	return 590*tiempoIda, 590*tiempoVuelta
}

func CompareInt(a int, b int) int{
	//returns 1 if a > b 0 otherwise
	if a>b{
		return 1
	} else {
		return 0
	}
}

func CreateReservation(c *gin.Context) {

	var reserva models.Reservation

	// Bind the request body to the reserva variable
	if err := c.BindJSON(&reserva); err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
	} else {

		fmt.Println("Generating PNR")
		reserva.PNR = GenerateNewPNR()

		fmt.Println("Calculating Balances")
		for i := 0; i<len(reserva.Passengers); i++{
			reserva.Passengers[i].Balances.AncillariesIda, reserva.Passengers[i].Balances.AncillariesVuelta = SumAncillaries(reserva.Passengers[i].Ancillaries)
			reserva.Passengers[i].Balances.VueloIda, reserva.Passengers[i].Balances.VueloVuelta = SumVuelos(c, reserva.Vuelos)
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
			reserva.Passengers[i].Balances.VueloIda, reserva.Passengers[i].Balances.VueloVuelta = SumVuelos(c, reserva.Vuelos)
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
