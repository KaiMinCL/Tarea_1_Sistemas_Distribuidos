package main

import (
	"bd_aerolinea/controllers"
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

	"github.com/joho/godotenv"
)

var (
	AncillaryPrice = map[string]int{"BGH": 10000, "BGR": 30000, "STDF": 5000, "PAXS": 2000, "PTCR": 40000, "AVIH": 40000, "SPML": 35000, "LNGE": 15000, "WIFI": 20000}
	meses          = []string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}
)

func waitAnimation() {
	fmt.Print("\nCargando")

	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		fmt.Print(".")
	}

	fmt.Println()
}

func clearScreen() {
	cmd := &exec.Cmd{}
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		SERVER = os.Getenv("SERVER")
		PORT   = os.Getenv("PORT")

		URL = "http://" + SERVER + ":" + PORT + "/api"

		option    int
		subOption int
	)

	for true {

		fmt.Println("Menu")
		time.Sleep(25 * time.Millisecond)

		fmt.Println("1. Gestionar reserva")
		time.Sleep(25 * time.Millisecond)

		fmt.Println("2. Obtener estadisticas")
		time.Sleep(25 * time.Millisecond)

		fmt.Println("3. Salir")
		time.Sleep(25 * time.Millisecond)

		fmt.Print("Ingrese una opcion: ")
		fmt.Scan(&option)

		switch option {
		case 1:
			clearScreen()

			fmt.Println("Submenu:")
			time.Sleep(25 * time.Millisecond)

			fmt.Println("1. Crear reserva")
			time.Sleep(25 * time.Millisecond)

			fmt.Println("2. Obtener reserva")
			time.Sleep(25 * time.Millisecond)

			fmt.Println("3. Modificar reserva")
			time.Sleep(25 * time.Millisecond)

			fmt.Println("4. Salir")
			time.Sleep(25 * time.Millisecond)

			fmt.Print("Ingrese una opcion: ")
			fmt.Scan(&subOption)

			switch subOption {
			case 1:
				clearScreen()
				// ALL OF THIS IS FOR CREATING A RESERVATION
				// It remains the following functionalities:
				// - check if there are seats left on the plane before selling them CHECK
				// - Remove capacity for every passenger sold on the flights CHECK
				// - Remove capacity for every ancillary sold on the flight CHECK
				// - Update mongoDB with the updated flight information CHECK

				var (
					fechaIda          string
					fechaRegreso      string
					origen            string
					destino           string
					cantidadPasajeros int
					Pasajeros         []models.Passenger
					Reserva           models.Reservation
					vueloVuelta       models.Flight
				)
				fmt.Println("Crear reserva")
				time.Sleep(25 * time.Millisecond)

				fmt.Print("Ingrese la fecha de ida: ")
				fmt.Scan(&fechaIda)
				time.Sleep(25 * time.Millisecond)

				fmt.Print("Ingrese la fecha de regreso: ")
				fmt.Scan(&fechaRegreso)
				time.Sleep(25 * time.Millisecond)

				fmt.Print("Ingrese aeropuerto de origen: ")
				fmt.Scan(&origen)
				time.Sleep(25 * time.Millisecond)

				fmt.Print("Ingrese aeropuerto de destino: ")
				fmt.Scan(&destino)
				time.Sleep(25 * time.Millisecond)

				fmt.Print("Ingrese la cantidad de pasajeros: ")
				fmt.Scan(&cantidadPasajeros)
				time.Sleep(25 * time.Millisecond)

				fmt.Println("")

				waitAnimation()
				clearScreen()

				resp, err := http.Get(URL + "/vuelo?origen=" + origen + "&destino=" + destino + "&fecha=" + fechaIda)

				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()

				// Read the response body
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("error:", err)
					break
				}

				var vuelos []models.Flight

				// parse the JSON response into the vuelos slice
				err = json.Unmarshal(body, &vuelos)
				if err != nil {
					fmt.Println("error:", err)
					break
				}

				// Checks if there's seats on the flight
				for i, vuelo := range vuelos {

					if vuelo.Avion.StockDePasajeros == 0 {
						vuelos = append(vuelos[:i], vuelos[i+1:]...)
						i -= 1
					}
				}

				// Check if the vuelos slice is empty
				if len(vuelos) == 0 {
					fmt.Println("No hay vuelos disponibles para la fecha de ida")
					break
				}

				// print out the vuelos slice to verify that it was parsed correctly

				fmt.Println("Vuelos disponibles:")
				fmt.Println("Ida:")
				for i := 0; i < len(vuelos); i++ {
					time.Sleep(50 * time.Millisecond)

					// Parse the start and end times

					horaSalida, _ := time.Parse("15:04", vuelos[i].HoraSalida)
					horaLlegada, _ := time.Parse("15:04", vuelos[i].HoraLlegada)

					if horaLlegada.Before(horaSalida) {
						horaLlegada = horaLlegada.Add(24 * time.Hour)
					}
					minutosVuelo := int(horaLlegada.Sub(horaSalida).Minutes())

					precioVuelo := 590 * minutosVuelo
					fmt.Print("\t", i+1)
					fmt.Print(". " + vuelos[i].NumeroVuelo + " " + vuelos[i].HoraSalida + " - " + vuelos[i].HoraLlegada + " $")
					fmt.Print(precioVuelo, "\n")
				}

				var opcionIda int

				fmt.Print("Ingrese una Opción: ")
				fmt.Scan(&opcionIda)

				vueloIda := vuelos[opcionIda-1]

				vueloIda.Avion.StockDePasajeros -= 1

				vueloReserva := models.ReservationFlight{
					NumeroVuelo: vueloIda.NumeroVuelo,
					Origen:      origen,
					Destino:     destino,
					HoraSalida:  vueloIda.HoraSalida,
					HoraLlegada: vueloIda.HoraLlegada,
					Fecha:       fechaIda,
				}
				Reserva.Vuelos = append(Reserva.Vuelos, vueloReserva)

				if fechaRegreso != "no" {
					resp, err := http.Get(URL + "/vuelo?origen=" + destino + "&destino=" + origen + "&fecha=" + fechaRegreso)

					if err != nil {
						log.Fatal(err)
						break
					}
					defer resp.Body.Close()

					// Read the response body
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Println("error:", err)
						break
					}

					var vuelos []models.Flight

					// parse the JSON response into the vuelos slice
					err = json.Unmarshal(body, &vuelos)
					if err != nil {
						fmt.Println("error:", err)
						break
					}

					// Checks if there's seats on the flight
					for i, vuelo := range vuelos {

						if vuelo.Avion.StockDePasajeros == 0 {
							vuelos = append(vuelos[:i], vuelos[i+1:]...)
							i -= 1
						}
					}

					// Check if the vuelos slice is empty
					if len(vuelos) == 0 {
						fmt.Println("No hay vuelos disponibles para la fecha de vuelta")
						break
					}

					fmt.Println("Vuelta:")

					var precioVuelo int

					// print out the vuelos slice to verify that it was parsed correctly
					for i := 0; i < len(vuelos); i++ {

						time.Sleep(50 * time.Millisecond)

						// Parse the start and end times

						horaSalida, _ := time.Parse("15:04", vuelos[i].HoraSalida)
						horaLlegada, _ := time.Parse("15:04", vuelos[i].HoraLlegada)

						if horaLlegada.Before(horaSalida) {
							horaLlegada = horaLlegada.Add(24 * time.Hour)
						}

						minutosVuelo := int(horaLlegada.Sub(horaSalida).Minutes())

						precioVuelo = 590 * minutosVuelo
						fmt.Print("\t", i+1)
						fmt.Print(". " + vuelos[i].NumeroVuelo + " " + vuelos[i].HoraSalida + " - " + vuelos[i].HoraLlegada + " $")
						fmt.Print(precioVuelo, "\n")
					}

					var opcionVuelta int

					fmt.Print("Ingrese una Opción: ")
					fmt.Scan(&opcionVuelta)

					vueloVuelta = vuelos[opcionVuelta-1]

					vueloVuelta.Avion.StockDePasajeros -= 1

					vueloReserva := models.ReservationFlight{
						NumeroVuelo: vueloVuelta.NumeroVuelo,
						Origen:      destino,
						Destino:     origen,
						HoraSalida:  vueloVuelta.HoraSalida,
						HoraLlegada: vueloVuelta.HoraLlegada,
						Fecha:       fechaRegreso,
					}

					Reserva.Vuelos = append(Reserva.Vuelos, vueloReserva)
				}

				for i := 0; i < cantidadPasajeros; i++ {

					clearScreen()

					fmt.Print("Pasajero ", i+1, " :\n")
					time.Sleep(25 * time.Millisecond)

					var Pasajero models.Passenger

					fmt.Print("Ingrese Nombre: ")
					fmt.Scan(&Pasajero.Name)
					time.Sleep(25 * time.Millisecond)

					fmt.Print("Ingrese Apellido: ")
					fmt.Scan(&Pasajero.Apellido)
					time.Sleep(25 * time.Millisecond)

					if i == 0 {
						Reserva.Apellido = Pasajero.Apellido
					}

					fmt.Print("Ingrese Edad: ")
					fmt.Scan(&Pasajero.Edad)
					time.Sleep(25 * time.Millisecond)

					fmt.Println("Ancillares Ida: ")

					for i := 0; i < len(vueloIda.Ancillaries); i++ {
						time.Sleep(50 * time.Millisecond)
						fmt.Print("\t", i+1)
						fmt.Print(". " + vueloIda.Ancillaries[i].Nombre + " | Stock: " + fmt.Sprint(vueloIda.Ancillaries[i].Stock) + " | Valor: $" + fmt.Sprint(AncillaryPrice[vueloIda.Ancillaries[i].SSR]) + "\n")
					}

					var seleccionAncillaries string

					fmt.Print("\nIngrese los Ancillaries (separados por comas): ")
					fmt.Scan(&seleccionAncillaries)

					stringSlice := strings.Split(seleccionAncillaries, ",")
					seleccionArray := make([]int, len(stringSlice))

					for i, v := range stringSlice {
						time.Sleep(50 * time.Millisecond)

						seleccionArray[i], err = strconv.Atoi(v)
						seleccionArray[i] -= 1

						var seleccionAncillary models.PassengerAncillary

						seleccionAncillary.SSR = vueloIda.Ancillaries[seleccionArray[i]].SSR

						if vueloIda.Ancillaries[seleccionArray[i]].Stock == 0 {
							fmt.Println("No existe stock para el Ancillary: " + vueloIda.Ancillaries[seleccionArray[i]].Nombre)
							break
						} else {
							for j := 0; j <= len(Pasajero.Ancillaries.Ida); j++ {

								vueloIda.Ancillaries[seleccionArray[i]].Stock -= 1

								if j == len(Pasajero.Ancillaries.Ida) {
									seleccionAncillary.Cantidad = 1
									break
								}

								if seleccionAncillary.SSR == Pasajero.Ancillaries.Ida[j].SSR {
									Pasajero.Ancillaries.Ida[j].Cantidad += 1
									break
								}

							}
						}
						Pasajero.Ancillaries.Ida = append(Pasajero.Ancillaries.Ida, seleccionAncillary)
					}

					if fechaRegreso != "no" {
						fmt.Println("Ancillares Vuelta: ")

						for i := 0; i < len(vueloVuelta.Ancillaries); i++ {
							time.Sleep(50 * time.Millisecond)

							fmt.Print("\t", i+1)
							fmt.Print(". " + vueloVuelta.Ancillaries[i].Nombre + " | Stock: " + fmt.Sprint(vueloVuelta.Ancillaries[i].Stock) + " | Valor: $" + fmt.Sprint(AncillaryPrice[vueloVuelta.Ancillaries[i].SSR]) + "\n")
						}

						var seleccionAncillaries string

						fmt.Print("\nIngrese los Ancillaries (separados por comas): ")
						fmt.Scan(&seleccionAncillaries)

						stringSlice := strings.Split(seleccionAncillaries, ",")
						seleccionArray := make([]int, len(stringSlice))

						for i, v := range stringSlice {

							time.Sleep(50 * time.Millisecond)

							seleccionArray[i], err = strconv.Atoi(v)
							seleccionArray[i] -= 1

							var seleccionAncillary models.PassengerAncillary

							seleccionAncillary.SSR = vueloVuelta.Ancillaries[seleccionArray[i]].SSR

							if vueloVuelta.Ancillaries[seleccionArray[i]].Stock == 0 {
								fmt.Println("No existe stock para el Ancillary: " + vueloVuelta.Ancillaries[seleccionArray[i]].Nombre)
								break
							} else {
								for j := 0; j <= len(Pasajero.Ancillaries.Vuelta); j++ {

									vueloVuelta.Ancillaries[seleccionArray[i]].Stock -= 1

									if j == len(Pasajero.Ancillaries.Vuelta) {
										seleccionAncillary.Cantidad = 1
										break
									}

									if seleccionAncillary.SSR == Pasajero.Ancillaries.Vuelta[j].SSR {
										Pasajero.Ancillaries.Vuelta[j].Cantidad += 1
										break
									}
								}
							}

							Pasajero.Ancillaries.Vuelta = append(Pasajero.Ancillaries.Vuelta, seleccionAncillary)
						}

					}

					Pasajeros = append(Pasajeros, Pasajero)
				}

				Reserva.Passengers = Pasajeros
				var PNRC models.PNRCapsule
				resp, err = http.Get(URL + "/generatepnr")

				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()

				body, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("error:", err)
					break
				}

				err = json.Unmarshal(body, &PNRC)
				if err != nil {
					fmt.Println(err)
					break
				}

				Reserva.PNR = PNRC.PNR

				for i := 0; i < len(Reserva.Passengers); i++ {
					Reserva.Passengers[i].Balances.AncillariesIda, Reserva.Passengers[i].Balances.AncillariesVuelta = controllers.SumAncillaries(Reserva.Passengers[i].Ancillaries)
					Reserva.Passengers[i].Balances.VueloIda, Reserva.Passengers[i].Balances.VueloVuelta = controllers.SumVuelos(Reserva.Vuelos)
				}

				JSONString, err := json.Marshal(Reserva)
				if err != nil {
					fmt.Println(err)
					return
				}

				resp, err = http.Post(URL+"/reserva", "application/json", bytes.NewBuffer(JSONString))

				if err != nil {
					log.Fatalf("Error realizando la reserva: %s", err.Error())
				}
				defer resp.Body.Close()

				resp, _ = http.Get(URL + "/reserva?pnr=" + Reserva.PNR + "&apellido=" + Reserva.Apellido)

				body, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("error:", err)
					return
				}

				// parse the JSON response into the vuelos slice
				err = json.Unmarshal(body, &Reserva)
				if err != nil {
					fmt.Println("error:", err)
					return
				}

				var update struct {
					StockDePasajeros int `json:"stock_de_pasajeros"`
				}

				update.StockDePasajeros = vueloIda.Avion.StockDePasajeros

				JSONString, err = json.Marshal(update)
				if err != nil {
					fmt.Println(err)
					return
				}

				req, err := http.NewRequest("PUT", URL+"/vuelo?numero_vuelo="+vueloIda.NumeroVuelo+"&origen="+vueloIda.Origen+"&destino="+vueloIda.Destino+"&fecha="+vueloIda.Fecha, bytes.NewBuffer(JSONString))

				if err != nil {
					// Handle error
					break
				}

				client := &http.Client{}
				resp, err = client.Do(req)

				if err != nil {
					// Handle error
					break
				}
				defer resp.Body.Close()

				if fechaRegreso != "no" {

					var update struct {
						StockDePasajeros int `json:"stock_de_pasajeros"`
					}

					update.StockDePasajeros = vueloVuelta.Avion.StockDePasajeros

					JSONString, err = json.Marshal(update)
					if err != nil {
						fmt.Println(err)
						return
					}

					req, err = http.NewRequest("PUT", URL+"/vuelo?numero_vuelo="+vueloVuelta.NumeroVuelo+"&origen="+vueloVuelta.Origen+"&destino="+vueloVuelta.Destino+"&fecha="+vueloVuelta.Fecha, bytes.NewBuffer(JSONString))

					if err != nil {
						// Handle error
						break
					}

					client = &http.Client{}
					resp, err = client.Do(req)

					if err != nil {
						// Handle error
						break
					}
					defer resp.Body.Close()
				}
				waitAnimation()
				clearScreen()

				fmt.Println("\nLa reserva fue generada con el PNR:", Reserva.PNR)

				var costoTotal int

				for _, pasajero := range Reserva.Passengers {
					costoTotal += pasajero.Balances.AncillariesIda + pasajero.Balances.VueloIda + pasajero.Balances.VueloVuelta + pasajero.Balances.VueloVuelta
				}

				fmt.Println("\nEl costo total de la reserva fue de: $" + fmt.Sprint(costoTotal))

				fmt.Print("\nPresione cualquier tecla para volver al menu principal...")
				bufio.NewReader(os.Stdin).ReadString('\n')

				var input string
				fmt.Scanln(&input)

			case 2:

				// ALL OF THIS IS FOR VIEWING A RESERVATION

				var (
					reserva models.Reservation
				)

				fmt.Print("Ingrese el PNR: ")
				fmt.Scan(&reserva.PNR)
				fmt.Print("Ingrese el Apellido: ")
				fmt.Scan(&reserva.Apellido)

				resp, err := http.Get(URL + "/reserva?pnr=" + reserva.PNR + "&apellido=" + reserva.Apellido)

				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()

				// Read the response body
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("error:", err)
					break
				}

				// parse the JSON response into the vuelos slice
				err = json.Unmarshal(body, &reserva)
				if err != nil {
					fmt.Println("No se ha podido encontrar la Reserva.")
					break
				}

				waitAnimation()
				clearScreen()

				fmt.Println("Reserva: ")

				fmt.Println("Vuelo Ida: ")

				fmt.Println("\tVuelo Nro: " + reserva.Vuelos[0].NumeroVuelo + " | Horario: " + reserva.Vuelos[0].HoraSalida + " - " + reserva.Vuelos[0].HoraLlegada)

				if len(reserva.Vuelos) != 1 {
					fmt.Println("Vuelo vuelta: ")
					fmt.Println("\tVuelo Nro: " + reserva.Vuelos[1].NumeroVuelo + " | Horario: " + reserva.Vuelos[1].HoraSalida + " - " + reserva.Vuelos[1].HoraLlegada)
				}

				fmt.Println("Pasajeros: ")
				for i := 0; i < len(reserva.Passengers); i++ {

					fmt.Println("\t", reserva.Passengers[i].Name, reserva.Passengers[i].Edad)

					fmt.Print("\t Ancillares ida:")

					for j := 0; j < len(reserva.Passengers[i].Ancillaries.Ida); j++ {
						fmt.Print(" " + reserva.Passengers[i].Ancillaries.Ida[j].SSR)
					}

					if len(reserva.Vuelos) != 1 {

						fmt.Print("\n\t Ancillares vuelta:")
						for j := 0; j < len(reserva.Passengers[i].Ancillaries.Vuelta); j++ {
							fmt.Print(" " + reserva.Passengers[i].Ancillaries.Vuelta[j].SSR)
						}

					}

					fmt.Println("")
				}
				fmt.Print("\nPresione cualquier tecla para volver al menu principal...")
				bufio.NewReader(os.Stdin).ReadString('\n')

				var input string
				fmt.Scanln(&input)

			case 3:

				// ALL OF THIS IS FOR MODIFYING A RESERVATION

				var (
					reserva models.Reservation
				)

				fmt.Print("Ingrese el PNR: ")
				fmt.Scan(&reserva.PNR)

				fmt.Print("Ingrese el Apellido: ")
				fmt.Scan(&reserva.Apellido)

				resp, err := http.Get(URL + "/reserva?pnr=" + reserva.PNR + "&apellido=" + reserva.Apellido)

				waitAnimation()
				clearScreen()

				if err != nil {
					log.Fatal(err)
					break
				}
				defer resp.Body.Close()

				// Read the response body
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("error:", err)
					break
				}

				// parse the JSON response into the vuelos slice
				err = json.Unmarshal(body, &reserva)
				if err != nil {
					fmt.Println("error:", err)
					break
				}

				fmt.Println("Opciones: ")
				time.Sleep(25 * time.Millisecond)

				fmt.Println("1. Cambiar fecha de vuelo")
				time.Sleep(25 * time.Millisecond)

				fmt.Println("2. Agregar Ancillaries")
				time.Sleep(25 * time.Millisecond)

				fmt.Println("3. Salir")
				time.Sleep(25 * time.Millisecond)

				var modifyOption int

				fmt.Print("Ingrese una opción: ")
				fmt.Scan(&modifyOption)
				time.Sleep(25 * time.Millisecond)

				switch modifyOption {
				case 1:

					// ALL OF THIS IS FOR MODIFYING THE FLIGHT
					// It remains the following functionalities:
					// - check if there are seats left on the NEW PLANE before selling them CHECK
					// - Remove capacity for every passenger sold on the flights CHECK
					// - ADD CAPACITY to the plane canceled CHECK
					// - Remove stock to every ancillary sold on the NEW FLIGHT CHECK
					// - ADD stock to every ancillary sold on the NEW FLIGHT CHECK
					// - Update mongoDB with the updated flight information CHECK

					waitAnimation()
					clearScreen()
					fmt.Println("Vuelos: ")
					time.Sleep(50 * time.Millisecond)

					fmt.Println("\t1. Ida: " + reserva.Vuelos[0].NumeroVuelo + " " + reserva.Vuelos[0].HoraSalida + " - " + reserva.Vuelos[0].HoraLlegada)

					if len(reserva.Vuelos) == 2 {
						time.Sleep(50 * time.Millisecond)
						fmt.Println("\t2. Vuelta: " + reserva.Vuelos[1].NumeroVuelo + " " + reserva.Vuelos[1].HoraSalida + " - " + reserva.Vuelos[1].HoraLlegada)
					}

					var flightReserved int
					var fecha string

					fmt.Print("Ingrese una opción: ")
					fmt.Scan(&flightReserved)
					time.Sleep(25 * time.Millisecond)

					flightReserved -= 1

					fmt.Print("Ingrese nueva fecha: ")
					fmt.Scan(&fecha)
					time.Sleep(25 * time.Millisecond)

					resp, err := http.Get(URL + "/vuelo?origen=" + reserva.Vuelos[flightReserved].Origen + "&destino=" + reserva.Vuelos[flightReserved].Destino + "&fecha=" + fecha)

					if err != nil {
						log.Fatal(err)
					}
					defer resp.Body.Close()

					// Read the response body
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Println("error:", err)
						break
					}

					waitAnimation()
					clearScreen()

					var vuelos []models.Flight

					// parse the JSON response into the vuelos slice
					err = json.Unmarshal(body, &vuelos)
					if err != nil {
						fmt.Println("error:", err)
						break
					}

					// Checks if there's seats on the flights
					for i, vuelo := range vuelos {

						if vuelo.Avion.StockDePasajeros < len(reserva.Passengers) || (vuelo.NumeroVuelo == reserva.Vuelos[flightReserved].NumeroVuelo && vuelo.Fecha == reserva.Vuelos[flightReserved].Fecha) {
							vuelos = append(vuelos[:i], vuelos[i+1:]...)
							i -= 1
						}

					}

					// Check if the vuelos slice is empty
					if len(vuelos) == 0 {
						fmt.Println("No hay vuelos disponibles para la fecha indicada")
						break
					}

					var precioVuelo int
					// print out the vuelos slice to verify that it was parsed correctly
					for i := 0; i < len(vuelos); i++ {
						time.Sleep(50 * time.Millisecond)

						// Parse the start and end times

						horaSalida, _ := time.Parse("15:04", vuelos[i].HoraSalida)
						horaLlegada, _ := time.Parse("15:04", vuelos[i].HoraLlegada)

						if horaLlegada.Before(horaSalida) {
							horaLlegada = horaLlegada.Add(24 * time.Hour)
						}
						minutosVuelo := int(horaLlegada.Sub(horaSalida).Minutes())

						precioVuelo = 590*minutosVuelo + 20000
						fmt.Print(i + 1)
						fmt.Print(". " + vuelos[i].NumeroVuelo + " " + vuelos[i].HoraSalida + " - " + vuelos[i].HoraLlegada + " $")
						fmt.Print(precioVuelo, "\n")
					}

					var flightOption int

					fmt.Print("Ingrese una Opción: ")
					fmt.Scan(&flightOption)
					time.Sleep(25 * time.Millisecond)

					nuevoVuelo := vuelos[flightOption-1]

					resp, err = http.Get(URL + "/vuelo?origen=" + reserva.Vuelos[flightReserved].Origen + "&destino=" + reserva.Vuelos[flightReserved].Destino + "&fecha=" + reserva.Vuelos[flightReserved].Fecha)

					if err != nil {
						log.Fatal(err)
					}
					defer resp.Body.Close()

					// Read the response body
					body, err = ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Println("error:", err)
						break
					}

					// parse the JSON response into the vuelos slice
					err = json.Unmarshal(body, &vuelos)
					if err != nil {
						fmt.Println("error:", err)
						break
					}

					var vueloViejo models.Flight

					for _, vuelo := range vuelos {

						if vuelo.NumeroVuelo == reserva.Vuelos[flightReserved].NumeroVuelo {
							vueloViejo = vuelos[0]
						}

					}

					/*vueloViejo.Avion.StockDePasajeros += len(reserva.Passengers)
					nuevoVuelo.Avion.StockDePasajeros -= len(reserva.Passengers)

					for i := 0; i < len(reserva.Passengers); i++ {
						if flightReserved == 0 {
							for j := 0; j < len(reserva.Passengers[i].Ancillaries.Ida); j++ {
								ancillary := models.FlightAncillary{
									Nombre: "",
									Stock:  reserva.Passengers[i].Ancillaries.Ida[j].Cantidad,
									SSR:    reserva.Passengers[i].Ancillaries.Ida[j].SSR,
								}

								for _, ancillaryViejo := range vueloViejo.Ancillaries {
									if ancillaryViejo.SSR == ancillary.SSR {
										ancillaryViejo.Stock += ancillary.Stock
										for _, ancillaryNuevo := range nuevoVuelo.Ancillaries {
											var updated = 0
											if ancillaryNuevo.SSR == ancillary.SSR {
												if ancillaryNuevo.Stock > 0 {
													ancillaryNuevo.Stock -= ancillary.Stock
												}
												updated = 1
											}
											if updated == 0 {
												ancillary.Nombre = ancillaryViejo.Nombre
												ancillary.Stock = 0
												nuevoVuelo.Ancillaries = append(nuevoVuelo.Ancillaries, ancillary)
											}
										}
									}

								}

							}

						} else {
							for j := 0; j < len(reserva.Passengers[i].Ancillaries.Vuelta); j++ {
								ancillary := models.FlightAncillary{
									Nombre: "",
									Stock:  reserva.Passengers[i].Ancillaries.Vuelta[j].Cantidad,
									SSR:    reserva.Passengers[i].Ancillaries.Vuelta[j].SSR,
								}

								for _, ancillaryViejo := range vueloViejo.Ancillaries {
									if ancillaryViejo.SSR == ancillary.SSR {
										ancillaryViejo.Stock += ancillary.Stock
										for _, ancillaryNuevo := range nuevoVuelo.Ancillaries {
											var updated = 0
											if ancillaryNuevo.SSR == ancillary.SSR {
												if ancillaryNuevo.Stock > 0 {
													ancillaryNuevo.Stock -= ancillary.Stock
												}
												updated = 1
											}
											if updated == 0 {
												ancillary.Nombre = ancillaryViejo.Nombre
												ancillary.Stock = 0
												nuevoVuelo.Ancillaries = append(nuevoVuelo.Ancillaries, ancillary)
											}
										}
									}

								}
							}

						}
					}*/

					vueloReserva := models.ReservationFlight{
						NumeroVuelo: nuevoVuelo.NumeroVuelo,
						Origen:      nuevoVuelo.Origen,
						Destino:     nuevoVuelo.Destino,
						HoraSalida:  nuevoVuelo.HoraSalida,
						HoraLlegada: nuevoVuelo.HoraLlegada,
						Fecha:       fecha,
					}
					reserva.Vuelos[flightReserved] = vueloReserva
					for i := 0; i < len(reserva.Passengers); i++ {
						if flightReserved == 0 {
							reserva.Passengers[i].Balances.VueloIda = precioVuelo
						} else {
							reserva.Passengers[i].Balances.VueloVuelta = precioVuelo
						}
					}

					JSONString, err := json.Marshal(reserva)
					if err != nil {
						fmt.Println(err)
						return
					}

					req, err := http.NewRequest("PUT", URL+"/reserva?pnr="+reserva.PNR+"&apellido="+reserva.Apellido, bytes.NewBuffer(JSONString))

					waitAnimation()
					clearScreen()

					if err != nil {
						fmt.Println("Error:", err)
						break
					}

					client := &http.Client{}
					resp, err = client.Do(req)

					if err != nil {
						fmt.Println("Error:", err)
						break
					}
					defer resp.Body.Close()

					var update struct {
						StockDePasajeros int `json:"stock_de_pasajeros"`
					}

					update.StockDePasajeros = vueloViejo.Avion.StockDePasajeros

					JSONString, err = json.Marshal(update)
					if err != nil {
						fmt.Println("Error:", err)
						return
					}

					req, err = http.NewRequest("PUT", URL+"/vuelo?numero_vuelo="+vueloViejo.NumeroVuelo+"&origen="+vueloViejo.Origen+"&destino="+vueloViejo.Destino+"&fecha="+vueloViejo.Fecha, bytes.NewBuffer(JSONString))

					if err != nil {
						fmt.Println("Error:", err)
						break
					}

					client = &http.Client{}
					resp, err = client.Do(req)

					if err != nil {
						fmt.Println("Error:", err)
						break
					}
					defer resp.Body.Close()

					update.StockDePasajeros = nuevoVuelo.Avion.StockDePasajeros

					JSONString, err = json.Marshal(update)
					if err != nil {
						fmt.Println("Error:", err)
						return
					}

					req, err = http.NewRequest("PUT", URL+"/vuelo?numero_vuelo="+nuevoVuelo.NumeroVuelo+"&origen="+nuevoVuelo.Origen+"&destino="+nuevoVuelo.Destino+"&fecha="+nuevoVuelo.Fecha, bytes.NewBuffer(JSONString))

					if err != nil {
						fmt.Println("Error:", err)
						break
					}

					client = &http.Client{}
					resp, err = client.Do(req)

					if err != nil {
						fmt.Println("Error:", err)
						break
					}
					defer resp.Body.Close()

					fmt.Println("¡La reserva fue modificada exitosamente!")

					fmt.Print("\nPresione cualquier tecla para volver al menu principal...")
					bufio.NewReader(os.Stdin).ReadString('\n')

					var input string
					fmt.Scanln(&input)

				case 2:

					// ALL OF THIS IS FOR ADDING ANCILLARIES TO THE FLIGHT
					// It remains the following functionalities:
					// - Remove stock to every new ancillary sold CHECK
					// - Update mongoDB with the updated flight information

					waitAnimation()
					clearScreen()

					var flightOption int

					fmt.Println("Vuelos: ")

					time.Sleep(50 * time.Millisecond)
					fmt.Println("\t1. Ida: " + reserva.Vuelos[0].NumeroVuelo + " " + reserva.Vuelos[0].HoraSalida + " - " + reserva.Vuelos[0].HoraLlegada)

					if len(reserva.Vuelos) == 2 {
						time.Sleep(50 * time.Millisecond)
						fmt.Println("\t2. Vuelta: " + reserva.Vuelos[1].NumeroVuelo + " " + reserva.Vuelos[1].HoraSalida + " - " + reserva.Vuelos[1].HoraLlegada)
					}

					fmt.Print("Ingrese una opción: ")
					fmt.Scan(&flightOption)

					flightOption -= 1

					resp, err := http.Get(URL + "/vuelo?origen=" + reserva.Vuelos[flightOption].Origen + "&destino=" + reserva.Vuelos[flightOption].Destino + "&fecha=" + reserva.Vuelos[flightOption].Fecha)

					waitAnimation()
					clearScreen()

					if err != nil {
						log.Fatal(err)
						break
					}
					defer resp.Body.Close()

					// Read the response body
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Println("error:", err)
						break
					}

					var vuelos []models.Flight

					// parse the JSON response into the vuelos slice
					err = json.Unmarshal(body, &vuelos)
					if err != nil {
						fmt.Println("error:", err)
						break
					}

					fmt.Println("Ancillaries disponibles: ")

					for i := 0; i < len(vuelos[0].Ancillaries); i++ {
						time.Sleep(50 * time.Millisecond)
						fmt.Print("\t", i+1)
						fmt.Print(". " + vuelos[0].Ancillaries[i].Nombre + " $" + fmt.Sprint(AncillaryPrice[vuelos[0].Ancillaries[i].SSR]) + "\n")
					}

					fmt.Println("\nPasajeros: ")

					for i := 0; i < len(reserva.Passengers); i++ {
						time.Sleep(50 * time.Millisecond)

						fmt.Println("\t", reserva.Passengers[i].Name, reserva.Passengers[i].Edad)

						fmt.Print("\nIngrese los Ancillaries (separados por comas): ")

						var seleccionAncillaries string
						fmt.Scan(&seleccionAncillaries)

						stringSlice := strings.Split(seleccionAncillaries, ",")
						seleccionArray := make([]int, len(stringSlice))

						for j, v := range stringSlice {
							seleccionArray[j], err = strconv.Atoi(v)
							seleccionArray[j] -= 1

							var seleccionAncillary models.PassengerAncillary

							seleccionAncillary.SSR = vuelos[0].Ancillaries[seleccionArray[j]].SSR

							if vuelos[0].Ancillaries[seleccionArray[j]].Stock == 0 {
								fmt.Println("No existe stock para el Ancillary: " + vuelos[0].Ancillaries[seleccionArray[j]].Nombre)
								break
							} else {
								for k := 0; k <= len(reserva.Passengers[i].Ancillaries.Ida); k++ {

									vuelos[0].Ancillaries[seleccionArray[j]].Stock -= 1

									if flightOption == 0 {
										if k == len(reserva.Passengers[i].Ancillaries.Ida) {
											seleccionAncillary.Cantidad = 1
											reserva.Passengers[i].Ancillaries.Ida = append(reserva.Passengers[i].Ancillaries.Ida, seleccionAncillary)
											reserva.Passengers[i].Balances.AncillariesIda += AncillaryPrice[reserva.Passengers[i].Ancillaries.Ida[k].SSR]
											break
										}

										if seleccionAncillary.SSR == reserva.Passengers[i].Ancillaries.Ida[k].SSR {
											reserva.Passengers[i].Ancillaries.Ida[k].Cantidad += 1
											reserva.Passengers[i].Balances.AncillariesIda += AncillaryPrice[reserva.Passengers[i].Ancillaries.Ida[k].SSR]
											break
										}

									} else {
										if k == len(reserva.Passengers[i].Ancillaries.Vuelta) {
											seleccionAncillary.Cantidad = 1
											reserva.Passengers[i].Ancillaries.Vuelta = append(reserva.Passengers[i].Ancillaries.Vuelta, seleccionAncillary)
											reserva.Passengers[i].Balances.AncillariesVuelta += AncillaryPrice[reserva.Passengers[i].Ancillaries.Vuelta[k].SSR]
											break
										}

										if seleccionAncillary.SSR == reserva.Passengers[i].Ancillaries.Vuelta[k].SSR {
											reserva.Passengers[i].Ancillaries.Vuelta[k].Cantidad += 1
											reserva.Passengers[i].Balances.AncillariesVuelta += AncillaryPrice[reserva.Passengers[i].Ancillaries.Vuelta[k].SSR]
											break
										}

									}

								}
							}

						}

						reserva.Passengers[i].Balances.AncillariesIda, reserva.Passengers[i].Balances.AncillariesVuelta = controllers.SumAncillaries(reserva.Passengers[i].Ancillaries)
					}

					JSONString, err := json.Marshal(reserva)
					if err != nil {
						fmt.Println(err)
						break
					}

					req, err := http.NewRequest("PUT", URL+"/reserva?pnr="+reserva.PNR+"&apellido="+reserva.Apellido, bytes.NewBuffer(JSONString))

					waitAnimation()
					clearScreen()

					if err != nil {
						// Handle error
					}

					client := &http.Client{}
					resp, err = client.Do(req)

					if err != nil {
						// Handle error
					}
					defer resp.Body.Close()

					/*var update struct {
						StockDePasajeros int `json:"stock_de_pasajeros"`
					}

					update.StockDePasajeros = vuelos[0].Avion.StockDePasajeros

					JSONString, err = json.Marshal(update)
					if err != nil {
						fmt.Println(err)
						return
					}

					req, err = http.NewRequest("PUT", URL+"/vuelo?numero_vuelo="+vuelos[0].NumeroVuelo+"&origen="+vuelos[0].Origen+"&destino="+vuelos[0].Destino+"&fecha="+vuelos[0].Fecha, bytes.NewBuffer(JSONString))

					if err != nil {
						// Handle error
						break
					}

					client = &http.Client{}
					resp, err = client.Do(req)

					if err != nil {
						// Handle error
						break
					}
					defer resp.Body.Close()*/

					fmt.Println("¡La reserva fue modificada exitosamente!")

					fmt.Print("\nPresione cualquier tecla para volver al menu principal...")
					bufio.NewReader(os.Stdin).ReadString('\n')

					var input string
					fmt.Scanln(&input)

				case 3:
					break

				default:
					fmt.Println("Opcion invalida")
				}
			case 4:
				break
			default:
				fmt.Println("Opcion invalida")
			}
		case 2:
			waitAnimation()
			clearScreen()
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

			fmt.Println("\nRuta Mayor Ganancia:", data["ruta_mayor_ganancia"], "\n")
			time.Sleep(50 * time.Millisecond)

			fmt.Println("Ruta Menor Ganancia:", data["ruta_menor_ganancia"], "\n")

			fmt.Println("Ranking Ancillaries:")
			for _, ra := range data["ranking_ancillaries"].([]interface{}) {

				fmt.Println("\tNombre:", ra.(map[string]interface{})["nombre"])
				time.Sleep(50 * time.Millisecond)

				fmt.Println("\tSSR:", ra.(map[string]interface{})["ssr"])
				time.Sleep(50 * time.Millisecond)

				fmt.Println("\tGanancia:", ra.(map[string]interface{})["ganancia"], "\n")
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

		case 3:
			clearScreen()
			return
		default:
			fmt.Println("Opcion invalida")
		}
	}
	clearScreen()
}
