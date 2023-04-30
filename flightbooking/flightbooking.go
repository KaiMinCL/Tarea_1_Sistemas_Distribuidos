package flightbooking

import (
	"bd_aerolinea/models"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	meses          = []string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}
	AncillaryPrice = map[string]int{"BGH": 10000, "BGR": 30000, "STDF": 5000, "PAXS": 2000, "PTCR": 40000, "AVIH": 40000, "SPML": 35000, "LNGE": 15000, "WIFI": 20000}
)

func WaitAnimation() {
	fmt.Print("\nCargando")

	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		fmt.Print(".")
	}

	fmt.Println()
}

func ClearScreen() {
	cmd := &exec.Cmd{}
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func reverseArray(vuelos []models.ReservationFlight) []models.ReservationFlight {
	for i := 0; i < len(vuelos)/2; i++ {
		j := len(vuelos) - 1 - i
		vuelos[i], vuelos[j] = vuelos[j], vuelos[i]
	}
	return vuelos
}

func GetVuelos(URL, origen, destino, fecha string) ([]models.Flight, error) {

	var vuelos []models.Flight
	resp, err := http.Get(URL + "/vuelo?origen=" + origen + "&destino=" + destino + "&fecha=" + fecha)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	// parse the JSON response into the vuelos slice
	err = json.Unmarshal(body, &vuelos)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}
	return vuelos, nil
}

// checkFlightAvailability checks if there are seats available on the flight.
func CheckFlightAvailability(URL, origen, destino, fecha string) ([]models.Flight, error) {
	vuelos, err := GetVuelos(URL, origen, destino, fecha)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	// Remove flights that don't have seats available
	for i := len(vuelos) - 1; i >= 0; i-- {
		if vuelos[i].Avion.StockDePasajeros == 0 {
			vuelos = append(vuelos[:i], vuelos[i+1:]...)
		}
	}

	return vuelos, nil
}

func calculateFlightPrice(flight models.Flight, modify bool) int {
	// Parse the start and end times
	horaSalida, _ := time.Parse("15:04", flight.HoraSalida)
	horaLlegada, _ := time.Parse("15:04", flight.HoraLlegada)

	if horaLlegada.Before(horaSalida) {
		horaLlegada = horaLlegada.Add(24 * time.Hour)
	}

	minutosVuelo := int(horaLlegada.Sub(horaSalida).Minutes())

	precioVuelo := 590 * minutosVuelo

	if modify {
		precioVuelo += 20000
	}

	return precioVuelo
}

// displayFlights displays the available flights.
func DisplayFlights(vuelos []models.Flight, vueloType string, modify bool) {
	fmt.Println("Vuelos disponibles:")
	fmt.Println(vueloType)
	for i, vuelo := range vuelos {
		time.Sleep(50 * time.Millisecond)

		precioVuelo := calculateFlightPrice(vuelo, modify)

		fmt.Printf("\t%d. %s %s - %s $%d\n", i+1, vuelo.NumeroVuelo, vuelo.HoraSalida, vuelo.HoraLlegada, precioVuelo)
	}
}

func UpdateFlightStock(URL string, vuelo models.Flight) error {
	update := struct {
		StockDePasajeros int `json:"stock_de_pasajeros"`
	}{}

	update.StockDePasajeros = vuelo.Avion.StockDePasajeros

	JSONString, err := json.Marshal(update)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", URL+"/vuelo?numero_vuelo="+vuelo.NumeroVuelo+"&origen="+vuelo.Origen+"&destino="+vuelo.Destino+"&fecha="+vuelo.Fecha, bytes.NewBuffer(JSONString))
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func SelectFlight(vuelos []models.Flight, reserva *models.Reservation, origen, destino, fechaIda string, modify bool) models.Flight {
	// print out the vuelos slice to verify that it was parsed correctly

	var opcion int

	fmt.Print("Ingrese una Opción: ")
	fmt.Scan(&opcion)

	vuelo := vuelos[opcion-1]

	vueloReserva := models.ReservationFlight{
		NumeroVuelo: vuelo.NumeroVuelo,
		Origen:      origen,
		Destino:     destino,
		HoraSalida:  vuelo.HoraSalida,
		HoraLlegada: vuelo.HoraLlegada,
		Fecha:       fechaIda,
	}

	if modify == true {
		fechaNuevoVuelo, _ := time.Parse("dd/mm/yyyy", vueloReserva.Fecha)
		fechaVueloActual, _ := time.Parse("dd/mm/yyyy", reserva.Vuelos[0].Fecha)

		if fechaNuevoVuelo.After(fechaVueloActual) {
			reserva.Vuelos = []models.ReservationFlight{reserva.Vuelos[0], vueloReserva}
		} else {
			reserva.Vuelos = []models.ReservationFlight{vueloReserva, reserva.Vuelos[0]}
		}

	} else {
		reserva.Vuelos = append(reserva.Vuelos, vueloReserva)
	}

	return vuelo
}

func parseAncillaryIndices(input string, maxIndex int) ([]int, error) {
	indicesStr := strings.Split(input, ",")
	indices := make([]int, len(indicesStr))

	for i, indexStr := range indicesStr {
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			return nil, fmt.Errorf("invalid index: %s", indexStr)
		}

		if index < 1 || index > maxIndex {
			return nil, fmt.Errorf("index out of range: %d", index)
		}

		indices[i] = index - 1
	}

	return indices, nil
}

func displayAncillaries(ancillaries []models.FlightAncillary, flightType string) {
	fmt.Printf("\nAncillaries %s:\n", flightType)

	for i, ancillary := range ancillaries {
		time.Sleep(50 * time.Millisecond)
		fmt.Printf("\t%d. %s | Stock: %d | Valor: $%d\n", i+1, ancillary.Nombre, ancillary.Stock, AncillaryPrice[ancillary.SSR])
	}
}

func SelectAncillaries(flightType string, ancillaries []models.FlightAncillary) ([]models.PassengerAncillary, error) {
	displayAncillaries(ancillaries, flightType)

	var selectedAncillaries []models.PassengerAncillary

	fmt.Printf("\nIngrese los Ancillaries %s (separados por comas): ", flightType)
	var input string
	fmt.Scan(&input)

	ssrIndices, err := parseAncillaryIndices(input, len(ancillaries))
	if err != nil {
		return nil, err
	}

	for _, ssrIndex := range ssrIndices {
		ancillary := ancillaries[ssrIndex]
		inAncillaries := false
		if ancillary.Stock == 0 {
			return nil, fmt.Errorf("no hay stock para el ancillary: %s", ancillary.Nombre)
		}

		var selectedAncillary models.PassengerAncillary
		selectedAncillary.SSR = ancillary.SSR
		selectedAncillary.Cantidad = 1

		for _, sa := range selectedAncillaries {
			if sa.SSR == selectedAncillary.SSR {
				selectedAncillary.Cantidad += 1
				inAncillaries = true
				break
			}
		}
		if !inAncillaries {
			selectedAncillaries = append(selectedAncillaries, selectedAncillary)
		}
	}

	return selectedAncillaries, nil
}

func AddAncillaries(selectedAncillaries []models.PassengerAncillary, flightOption, i int, reserva *models.Reservation) {
	for _, newAncillary := range selectedAncillaries {
		inAncillaries := false
		newAncillary.Cantidad = 1
		if flightOption == 0 {
			for k := range reserva.Passengers[i].Ancillaries.Ida {
				if newAncillary.SSR == reserva.Passengers[i].Ancillaries.Ida[k].SSR {
					reserva.Passengers[i].Ancillaries.Ida[k].Cantidad += 1
					inAncillaries = true
					break
				}
			}
			if !inAncillaries {
				reserva.Passengers[i].Ancillaries.Ida = append(reserva.Passengers[i].Ancillaries.Ida, newAncillary)
			}
			reserva.Passengers[i].Balances.AncillariesIda += AncillaryPrice[newAncillary.SSR]

		} else {
			for k := range reserva.Passengers[i].Ancillaries.Vuelta {
				if newAncillary.SSR == reserva.Passengers[i].Ancillaries.Vuelta[k].SSR {
					reserva.Passengers[i].Ancillaries.Vuelta[k].Cantidad += 1
					inAncillaries = true
					break
				}
			}
			if !inAncillaries {
				reserva.Passengers[i].Ancillaries.Vuelta = append(reserva.Passengers[i].Ancillaries.Vuelta, newAncillary)
			}
			reserva.Passengers[i].Balances.AncillariesVuelta += AncillaryPrice[newAncillary.SSR]

		}
	}
}

func PassengerInformation(cantidadPasajeros int, vuelos []models.Flight) ([]models.Passenger, error) {
	var pasajeros []models.Passenger

	for i := 0; i < cantidadPasajeros; i++ {
		ClearScreen()
		fmt.Printf("Pasajero %d:\n", i+1)
		time.Sleep(25 * time.Millisecond)

		var pasajero models.Passenger

		fmt.Print("Ingrese Nombre: ")
		fmt.Scan(&pasajero.Name)
		time.Sleep(25 * time.Millisecond)

		fmt.Print("Ingrese Apellido: ")
		fmt.Scan(&pasajero.Apellido)
		time.Sleep(25 * time.Millisecond)

		fmt.Print("Ingrese Edad: ")
		fmt.Scan(&pasajero.Edad)
		time.Sleep(25 * time.Millisecond)

		idaAncillaries, err := SelectAncillaries("Ida", vuelos[0].Ancillaries)
		if err != nil {
			return nil, err
		}
		pasajero.Ancillaries.Ida = idaAncillaries

		if len(vuelos) == 2 {
			vueltaAncillaries, err := SelectAncillaries("Vuelta", vuelos[1].Ancillaries)
			if err != nil {
				return nil, err
			}
			pasajero.Ancillaries.Vuelta = vueltaAncillaries
		}

		pasajeros = append(pasajeros, pasajero)
	}

	return pasajeros, nil
}

func GeneratePNR(URL string) (string, error) {
	resp, err := http.Get(URL + "/generatepnr")
	if err != nil {
		return "", fmt.Errorf("failed to make HTTP GET request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var PNRC models.PNRCapsule

	err = json.Unmarshal(body, &PNRC)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON response: %v", err)
	}

	return PNRC.PNR, nil
}

func GetReserva(URL string, reserva *models.Reservation) error {
	resp, err := http.Get(URL + "/reserva?pnr=" + reserva.PNR + "&apellido=" + strings.Title(reserva.Apellido))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// parse the JSON response into the vuelos slice
	err = json.Unmarshal(body, &reserva)
	if err != nil {
		return err
	}

	return nil
}

func MakeReservation(URL string, reserva models.Reservation, modify bool) error {
	// Update passenger balances
	for i := 0; i < len(reserva.Passengers); i++ {
		reserva.Passengers[i].Balances.AncillariesIda, reserva.Passengers[i].Balances.AncillariesVuelta = models.SumAncillaries(reserva.Passengers[i].Ancillaries)
		reserva.Passengers[i].Balances.VueloIda, reserva.Passengers[i].Balances.VueloVuelta = models.SumVuelos(reserva.Vuelos)
	}

	// Make HTTP request to create the reservation
	JSONString, err := json.Marshal(reserva)
	if err != nil {
		return err
	}

	if modify == true {

		req, err := http.NewRequest("PUT", URL+"/reserva?pnr="+reserva.PNR+"&apellido="+reserva.Apellido, bytes.NewBuffer(JSONString))
		if err != nil {
			return fmt.Errorf("Error realizando la reserva: %s", err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	} else {
		resp, err := http.Post(URL+"/reserva", "application/json", bytes.NewBuffer(JSONString))
		if err != nil {
			return fmt.Errorf("Error realizando la reserva: %s", err.Error())
		}
		defer resp.Body.Close()

	}

	return nil
}

func DisplayStatistics(URL string) {
	WaitAnimation()
	ClearScreen()
	resp, err := http.Get(URL + "/estadisticas")

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	var data map[string]interface{}

	// parse the JSON response into the vuelos slice
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		panic(err)
	}

	fmt.Println("Estadisticas:")
	time.Sleep(50 * time.Millisecond)

	fmt.Println("\nRuta Mayor Ganancia:", data["ruta_mayor_ganancia"])
	time.Sleep(50 * time.Millisecond)

	fmt.Println("Ruta Menor Ganancia:", data["ruta_menor_ganancia"])

	fmt.Println("Ranking Ancillaries:")
	for _, ra := range data["ranking_ancillaries"].([]interface{}) {

		fmt.Println("\tNombre:", ra.(map[string]interface{})["nombre"])
		time.Sleep(50 * time.Millisecond)

		fmt.Println("\tSSR:", ra.(map[string]interface{})["ssr"])
		time.Sleep(50 * time.Millisecond)

		fmt.Println("\tGanancia:", ra.(map[string]interface{})["ganancia"])
	}
	fmt.Println("Promedio Pasajeros:")
	for _, mes := range meses {
		valor, ok := data["promedio_pasajeros"].(map[string]interface{})[mes]
		if !ok {
			continue
		}
		time.Sleep(50 * time.Millisecond)
		fmt.Println("\t", strings.Title(mes), ":", valor)
	}

	fmt.Print("\nPresione cualquier tecla para volver al menu principal...")
	bufio.NewReader(os.Stdin).ReadString('\n')

	var input string
	fmt.Scanln(&input)
}
