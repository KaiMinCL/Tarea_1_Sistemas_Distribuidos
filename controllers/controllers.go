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
	AncillaryName = map[string]string{"BGH":"Equipaje de mano", "BGR":"Equipaje de bodega", "STDF":"Asiento", "PAXS":"Embarque y Check In prioritario", "PTCR":"Mascota en cabina", "AVIH":"Mascota en bodega", "SPML":"Equipaje especial", "LNGE":"Acceso a Salon VIP", "WIFI":"Wi-Fi a bordo"}
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

func Max(data map[string]int) string{
	var max int = 0
	var maxIndex string
	for i, v := range(data){
		if v > max {
			maxIndex = i
			max = v
		}
	}
	return maxIndex
}

func Min(data map[string]int) string{
	var min int = 2000000000
	var minIndex string
	for i, v := range(data){
		if v < min {
			minIndex = i
			min = v
		}
	}
	return minIndex
}



func GetStatistics(c *gin.Context){
	var(
		balancesVuelo = make(map[string]int)
		balancesAncillaries = make(map[string]int)
		stats models.Statistics
	)
	reservas, err := models.GetAllReservations()

	if err != nil{
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	for i := 0; i<len(reservas); i++{
		for j := 0; j<len(reservas[i].Passengers); j++{
			balancesVuelo[reservas[i].Vuelos[0].NumeroVuelo] += reservas[i].Passengers[j].Balances.VueloIda
			if len(reservas[i].Vuelos) == 2{
				balancesVuelo[reservas[i].Vuelos[1].NumeroVuelo] += reservas[i].Passengers[j].Balances.VueloVuelta
			}

			for k := 0; k<len(reservas[i].Passengers[j].Ancillaries.Ida); k++{
				balancesAncillaries[reservas[i].Passengers[j].Ancillaries.Ida[k].SSR] += reservas[i].Passengers[j].Balances.AncillariesIda
			}
			for k := 0; k<len(reservas[i].Passengers[j].Ancillaries.Vuelta); k++{
				balancesAncillaries[reservas[i].Passengers[j].Ancillaries.Vuelta[k].SSR] += reservas[i].Passengers[j].Balances.AncillariesVuelta
			}
		}
	}
	stats.RutaMayorGanancia = Max(balancesVuelo)
	stats.RutaMenorGanancia = Min(balancesVuelo)
	for i, v := range(AncillaryName){

		var RankingAncillary models.AncillariesStatistics
		RankingAncillary.Nombre = v
		RankingAncillary.SSR = i
		RankingAncillary.Ganancia = balancesAncillaries[i]
		stats.RankingAncillaries = append(stats.RankingAncillaries, RankingAncillary)
	}

	layout := "02/01/2006"

	//var months = [12]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	var gananciaMes map[string]int = map[string]int{"January":0, "February":0, "March":0, "April":0, "May":0, "June":0, "July":0, "August":0, "September":0, "October":0, "November":0, "December":0}
	var pasajerosMes map[string]int = map[string]int{"January":0, "February":0, "March":0, "April":0, "May":0, "June":0, "July":0, "August":0, "September":0, "October":0, "November":0, "December":0}

	for i := 0; i<len(reservas); i++{
		dateIda := reservas[i].Vuelos[0].Fecha
		fmt.Println(dateIda)
		mesIda, _ := time.Parse(layout, dateIda)
		fmt.Println(fmt.Sprint(mesIda.Month()))

		pasajerosMes[fmt.Sprint(mesIda.Month())] += len(reservas[i].Passengers)
		for j := 0; j<len(reservas[i].Passengers); j++{
			gananciaMes[fmt.Sprint(mesIda.Month())] += reservas[i].Passengers[j].Balances.VueloIda
			gananciaMes[fmt.Sprint(mesIda.Month())] += reservas[i].Passengers[j].Balances.AncillariesIda
		}
		if len(reservas[i].Vuelos) == 2{
			dateVuelta := reservas[i].Vuelos[1].Fecha
			mesVuelta, _ := time.Parse(layout, dateVuelta)
			pasajerosMes[fmt.Sprint(mesVuelta.Month())] += len(reservas[i].Passengers)

			for j := 0; j<len(reservas[i].Passengers); j++{
				gananciaMes[fmt.Sprint(mesVuelta.Month())] += reservas[i].Passengers[j].Balances.VueloVuelta
				gananciaMes[fmt.Sprint(mesVuelta.Month())] += reservas[i].Passengers[j].Balances.AncillariesVuelta
			}
		}
	}
	fmt.Print(gananciaMes)
	fmt.Print(pasajerosMes)

	stats.PromedioPasajeros.Jan = GetAverage(gananciaMes["January"], pasajerosMes["January"])
	stats.PromedioPasajeros.Feb = GetAverage(gananciaMes["February"], pasajerosMes["February"])
	stats.PromedioPasajeros.Mar = GetAverage(gananciaMes["March"], pasajerosMes["March"])
	stats.PromedioPasajeros.Apr = GetAverage(gananciaMes["April"], pasajerosMes["April"])
	stats.PromedioPasajeros.May = GetAverage(gananciaMes["May"], pasajerosMes["May"])
	stats.PromedioPasajeros.Jun = GetAverage(gananciaMes["June"], pasajerosMes["June"])
	stats.PromedioPasajeros.Jul = GetAverage(gananciaMes["July"], pasajerosMes["July"])
	stats.PromedioPasajeros.Aug = GetAverage(gananciaMes["August"], pasajerosMes["August"])
	stats.PromedioPasajeros.Sep = GetAverage(gananciaMes["September"], pasajerosMes["September"])
	stats.PromedioPasajeros.Oct = GetAverage(gananciaMes["October"], pasajerosMes["October"])
	stats.PromedioPasajeros.Nov = GetAverage(gananciaMes["November"], pasajerosMes["November"])
	stats.PromedioPasajeros.Dec = GetAverage(gananciaMes["December"], pasajerosMes["December"])

	c.IndentedJSON(http.StatusOK, stats)
}


func GetAverage(ganancia int, passenger int) int{
	if (passenger == 0){
		return 0
	} else {
		return int(ganancia/passenger)
	}
}
