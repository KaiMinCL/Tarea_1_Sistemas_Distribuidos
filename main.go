package main

import (
	"bd_aerolinea/controllers"
	"bd_aerolinea/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var (
	AncillaryPrice = map[string]int{"BGH": 10000, "BGR": 30000, "STDF": 5000, "PAXS": 2000, "PTCR": 40000, "AVIH": 40000, "SPML": 35000, "LNGE": 15000, "WIFI": 20000}
)

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

	fmt.Println("connecting to", SERVER, ":", PORT)

	for {
		fmt.Println("Menu")
		fmt.Println("1. Gestionar reserva")
		fmt.Println("2. Obtener estadisticas")
		fmt.Println("3. Salir")
		fmt.Print("Ingrese una opcion: ")
		fmt.Scan(&option)
		switch option {
		case 1:
			fmt.Println("Submenu:")
			fmt.Println("1. Crear reserva")
			fmt.Println("2. Obtener reserva")
			fmt.Println("3. Modificar reserva")
			fmt.Println("4. Salir")
			fmt.Print("Ingrese una opcion: ")
			fmt.Scan(&subOption)
			switch subOption {
			case 1:

				// ALL OF THIS IS FOR CREATING A RESERVATION
				// It remains the following functionalities:
				// - check if there are seats left on the plane before selling them
				// - Remove capacity for every passenger sold on the flights
				// - Remove capacity for every ancillary sold on the flight
				// - Update mongoDB with the updated flight information
				// - Indifference to string caps for the Reservation "Apellido" (Lastname)

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
				fmt.Println("Ingrese la fecha de ida:")
				fmt.Scan(&fechaIda)
				fmt.Println("Ingrese la fecha de regreso:")
				fmt.Scan(&fechaRegreso)
				fmt.Println("Ingrese origen:")
				fmt.Scan(&origen)
				fmt.Println("Ingrese destino")
				fmt.Scan(&destino)
				fmt.Println("Cantidad de Pasajeros:")
				fmt.Scan(&cantidadPasajeros)
				fmt.Println("Vuelos disponibles:")
				fmt.Println("Ida:")

				resp, err := http.Get(URL + "/vuelo?origen=" + origen + "&destino=" + destino + "&fecha=" + fechaIda)

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

				var vuelos []models.Flight

				// parse the JSON response into the vuelos slice
				err = json.Unmarshal(body, &vuelos)
				if err != nil {
					fmt.Println("error:", err)
					return
				}

				// Check if the vuelos slice is empty
				if len(vuelos) == 0 {
					fmt.Println("No hay vuelos disponibles para la fecha de ida")
					return
				}

				// print out the vuelos slice to verify that it was parsed correctly
				for i := 0; i < len(vuelos); i++ {

					// Parse the start and end times

					horaSalida, _ := time.Parse("15:04", vuelos[i].HoraSalida)
					horaLlegada, _ := time.Parse("15:04", vuelos[i].HoraLlegada)

					if horaLlegada.Before(horaSalida) {
						horaLlegada = horaLlegada.Add(24 * time.Hour)
					}
					minutosVuelo := int(horaLlegada.Sub(horaSalida).Minutes())

					precioVuelo := 590 * minutosVuelo
					fmt.Print(i + 1)
					fmt.Print(". " + vuelos[i].NumeroVuelo + " " + vuelos[i].HoraSalida + " - " + vuelos[i].HoraLlegada + " $")
					fmt.Print(precioVuelo, "\n")
				}

				var opcionIda int
				fmt.Print("\nIngrese una Opción: ")
				fmt.Scan(&opcionIda)
				vueloIda := vuelos[opcionIda-1]

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
					}
					defer resp.Body.Close()

					// Read the response body
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Println("error:", err)
						return
					}

					var vuelos []models.Flight

					// parse the JSON response into the vuelos slice
					err = json.Unmarshal(body, &vuelos)
					if err != nil {
						fmt.Println("error:", err)
						return
					}

					// Check if the vuelos slice is empty
					if len(vuelos) == 0 {
						fmt.Println("No hay vuelos disponibles para la fecha de vuelta")
						return
					}

					// print out the vuelos slice to verify that it was parsed correctly
					for i := 0; i < len(vuelos); i++ {

						// Parse the start and end times

						horaSalida, _ := time.Parse("15:04", vuelos[i].HoraSalida)
						horaLlegada, _ := time.Parse("15:04", vuelos[i].HoraLlegada)

						if horaLlegada.Before(horaSalida) {
							horaLlegada = horaLlegada.Add(24 * time.Hour)
						}

						minutosVuelo := int(horaLlegada.Sub(horaSalida).Minutes())

						precioVuelo := 590 * minutosVuelo
						fmt.Print(i + 1)
						fmt.Print(". " + vuelos[i].NumeroVuelo + " " + vuelos[i].HoraSalida + " - " + vuelos[i].HoraLlegada + " $")
						fmt.Print(precioVuelo, "\n")
					}

					var opcionVuelta int
					fmt.Println("Ingrese una Opción: ")
					fmt.Scan(&opcionVuelta)
					vueloVuelta = vuelos[opcionVuelta-1]

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

					fmt.Print("\nPasajero ", i+1, " :\n")

					var Pasajero models.Passenger

					fmt.Println("Ingrese Nombre: ")
					fmt.Scan(&Pasajero.Name)
					fmt.Println("Ingrese Apellido: ")
					fmt.Scan(&Pasajero.Apellido)

					if i == 0 {
						Reserva.Apellido = Pasajero.Apellido
					}

					fmt.Println("Ingrese Edad: ")
					fmt.Scan(&Pasajero.Edad)

					fmt.Println("Ancillares Ida: ")

					for i := 0; i < len(vueloIda.Ancillaries); i++ {
						fmt.Print(i + 1)
						fmt.Print(". " + vueloIda.Ancillaries[i].Nombre + " $" + fmt.Sprint(AncillaryPrice[vueloIda.Ancillaries[i].SSR]) + "\n")
					}

					fmt.Print("\nIngrese los Ancillaries (separados por comas): ")
					var seleccionAncillaries string
					fmt.Scan(&seleccionAncillaries)

					stringSlice := strings.Split(seleccionAncillaries, ",")
					seleccionArray := make([]int, len(stringSlice))

					for i, v := range stringSlice {
						seleccionArray[i], err = strconv.Atoi(v)
						seleccionArray[i] -= 1

						var seleccionAncillary models.PassengerAncillary

						seleccionAncillary.SSR = vueloIda.Ancillaries[seleccionArray[i]].SSR

						for j := 0; j <= len(Pasajero.Ancillaries.Ida); j++ {
							if j == len(Pasajero.Ancillaries.Ida) {
								seleccionAncillary.Cantidad = 1
								break
							}

							if seleccionAncillary.SSR == Pasajero.Ancillaries.Ida[j].SSR {
								Pasajero.Ancillaries.Ida[j].Cantidad += 1
								break
							}

						}
						Pasajero.Ancillaries.Ida = append(Pasajero.Ancillaries.Ida, seleccionAncillary)
					}
					if fechaRegreso != "no" {
						fmt.Println("Ancillares Vuelta: ")

						for i := 0; i < len(vueloVuelta.Ancillaries); i++ {
							fmt.Print(i + 1)
							fmt.Print(". " + vueloVuelta.Ancillaries[i].Nombre + " $" + fmt.Sprint(AncillaryPrice[vueloIda.Ancillaries[i].SSR]) + "\n")
						}

						fmt.Print("\nIngrese los Ancillaries (separados por comas): ")
						var seleccionAncillaries string
						fmt.Scan(&seleccionAncillaries)

						stringSlice := strings.Split(seleccionAncillaries, ",")
						seleccionArray := make([]int, len(stringSlice))

						for i, v := range stringSlice {
							seleccionArray[i], err = strconv.Atoi(v)
							seleccionArray[i] -= 1

							var seleccionAncillary models.PassengerAncillary

							seleccionAncillary.SSR = vueloVuelta.Ancillaries[seleccionArray[i]].SSR

							for j := 0; j <= len(Pasajero.Ancillaries.Vuelta); j++ {
								if j == len(Pasajero.Ancillaries.Vuelta) {
									seleccionAncillary.Cantidad = 1
									break
								}

								if seleccionAncillary.SSR == Pasajero.Ancillaries.Vuelta[j].SSR {
									Pasajero.Ancillaries.Vuelta[j].Cantidad += 1
									break
								}

							}
							Pasajero.Ancillaries.Vuelta = append(Pasajero.Ancillaries.Vuelta, seleccionAncillary)
						}

					}

					Pasajeros = append(Pasajeros, Pasajero)
				}

				Reserva.Passengers = Pasajeros
				Reserva.PNR = controllers.GenerateNewPNR()

				JSONString, err := json.Marshal(Reserva)
				if err != nil {
					fmt.Println(err)
					return
				}
				http.Post(URL+"/reserva", "application/json", bytes.NewBuffer(JSONString))

			case 2:

				fmt.Println("Obtener reserva")
				var (
					reserva models.Reservation
				)

				fmt.Println("Ingrese el PNR: ")
				fmt.Scan(&reserva.PNR)
				fmt.Println("Ingrese el Apellido: ")
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
					return
				}

				// parse the JSON response into the vuelos slice
				err = json.Unmarshal(body, &reserva)
				if err != nil {
					fmt.Println("error:", err)
					return
				}

				fmt.Println("Ida: ")

				fmt.Println(reserva.Vuelos[0].NumeroVuelo + " " + reserva.Vuelos[0].HoraSalida + " - " + reserva.Vuelos[0].HoraLlegada)

				fmt.Println("Vuelta: ")
				fmt.Println(reserva.Vuelos[1].NumeroVuelo + " " + reserva.Vuelos[1].HoraSalida + " - " + reserva.Vuelos[1].HoraLlegada)

				fmt.Println("Pasajeros: ")
				for i := 0; i < len(reserva.Passengers); i++ {
					fmt.Println(reserva.Passengers[i].Name, reserva.Passengers[i].Edad)
					fmt.Println("Ancillares ida:")
					for j := 0; j < len(reserva.Passengers[i].Ancillaries.Ida); j++ {
						fmt.Print(" " + reserva.Passengers[i].Ancillaries.Ida[j].SSR)
					}
					fmt.Println("\nAncillares vuelta:")
					for j := 0; j < len(reserva.Passengers[i].Ancillaries.Vuelta); j++ {
						fmt.Print(" " + reserva.Passengers[i].Ancillaries.Vuelta[j].SSR)
					}
				}
				fmt.Println("")

			case 3:
				fmt.Println("Modificar reserva")
				var (
					reserva models.Reservation
				)

				fmt.Println("Ingrese el PNR: ")
				fmt.Scan(&reserva.PNR)
				fmt.Println("Ingrese el Apellido: ")
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
					return
				}

				// parse the JSON response into the vuelos slice
				err = json.Unmarshal(body, &reserva)
				if err != nil {
					fmt.Println("error:", err)
					return
				}
				fmt.Println("Opciones: ")
				fmt.Println("1. Cambiar fecha de vuelo")
				fmt.Println("2. Agregar Ancillaries")
				fmt.Println("3. Salir")

				fmt.Println("Ingrese una opción: ")

				var modifyOption int
				fmt.Scan(&modifyOption)
				switch modifyOption {
				case 1:
					fmt.Println("Vuelos: ")

					fmt.Println("1. Ida: " + reserva.Vuelos[0].NumeroVuelo + " " + reserva.Vuelos[0].HoraSalida + " - " + reserva.Vuelos[0].HoraLlegada)
					fmt.Println("2. Vuelta: " + reserva.Vuelos[1].NumeroVuelo + " " + reserva.Vuelos[1].HoraSalida + " - " + reserva.Vuelos[1].HoraLlegada)

					var flightReserved int
					var fecha string
					fmt.Println("Ingrese una opción: ")
					fmt.Scan(&flightReserved)
					flightReserved -= 1

					fmt.Println("Ingrese nueva fecha: ")
					fmt.Scan(&fecha)
					resp, err := http.Get(URL + "/vuelo?origen=" + reserva.Vuelos[flightReserved].Origen + "&destino=" + reserva.Vuelos[flightReserved].Destino + "&fecha=" + fecha)

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

					var vuelos []models.Flight

					// parse the JSON response into the vuelos slice
					err = json.Unmarshal(body, &vuelos)
					if err != nil {
						fmt.Println("error:", err)
						return
					}

					// Check if the vuelos slice is empty
					if len(vuelos) == 0 {
						fmt.Println("No hay vuelos disponibles para la fecha de ida")
						return
					}

					// print out the vuelos slice to verify that it was parsed correctly
					for i := 0; i < len(vuelos); i++ {

						// Parse the start and end times

						horaSalida, _ := time.Parse("15:04", vuelos[i].HoraSalida)
						horaLlegada, _ := time.Parse("15:04", vuelos[i].HoraLlegada)

						if horaLlegada.Before(horaSalida) {
							horaLlegada = horaLlegada.Add(24 * time.Hour)
						}
						minutosVuelo := int(horaLlegada.Sub(horaSalida).Minutes())

						precioVuelo := 590 * minutosVuelo
						fmt.Print(i + 1)
						fmt.Print(". " + vuelos[i].NumeroVuelo + " " + vuelos[i].HoraSalida + " - " + vuelos[i].HoraLlegada + " $")
						fmt.Print(precioVuelo, "\n")
					}
					fmt.Print("\nIngrese una Opción: ")
					var flightOption int
					fmt.Scan(&flightOption)
					flightOption -= 1

					vuelo := vuelos[flightOption]

					vueloReserva := models.ReservationFlight{
						NumeroVuelo: vuelo.NumeroVuelo,
						Origen:      vuelo.Origen,
						Destino:     vuelo.Destino,
						HoraSalida:  vuelo.HoraSalida,
						HoraLlegada: vuelo.HoraLlegada,
						Fecha:       fecha,
					}
					reserva.Vuelos[flightReserved] = vueloReserva

					JSONString, err := json.Marshal(reserva)
					if err != nil {
						fmt.Println(err)
						return
					}

					req, err := http.NewRequest("PUT", URL+"/reserva?pnr="+reserva.PNR+"&apellido="+reserva.Apellido, bytes.NewBuffer(JSONString))

					if err != nil {
						// Handle error
					}

					client := &http.Client{}
					resp, err = client.Do(req)

					if err != nil {
						// Handle error
					}
					defer resp.Body.Close()
				case 2:
					fmt.Println("Vuelos: ")

					fmt.Println("1. Ida: " + reserva.Vuelos[0].NumeroVuelo + " " + reserva.Vuelos[0].HoraSalida + " - " + reserva.Vuelos[0].HoraLlegada)
					fmt.Println("2. Vuelta: " + reserva.Vuelos[1].NumeroVuelo + " " + reserva.Vuelos[1].HoraSalida + " - " + reserva.Vuelos[1].HoraLlegada)

					var flightOption int
					fmt.Println("Ingrese una opción: ")
					fmt.Scan(&flightOption)
					flightOption -= 1
					resp, err := http.Get(URL + "/vuelo?origen=" + reserva.Vuelos[flightOption].Origen + "&destino=" + reserva.Vuelos[flightOption].Destino + "&fecha=" + reserva.Vuelos[flightOption].Fecha)

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

					var vuelos []models.Flight

					// parse the JSON response into the vuelos slice
					err = json.Unmarshal(body, &vuelos)
					if err != nil {
						fmt.Println("error:", err)
						return
					}

					fmt.Println("Ancillaries disponibles: ")

					for i := 0; i < len(vuelos[0].Ancillaries); i++ {
						fmt.Print(i + 1)
						fmt.Print(". " + vuelos[0].Ancillaries[i].Nombre + " $" + fmt.Sprint(AncillaryPrice[vuelos[0].Ancillaries[i].SSR]) + "\n")
					}

					fmt.Println("Pasajeros: ")
					for i := 0; i < len(reserva.Passengers); i++ {

						fmt.Println(reserva.Passengers[i].Name, reserva.Passengers[i].Edad)
						fmt.Println("Ancillares ida:")

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

							for k := 0; k <= len(reserva.Passengers[i].Ancillaries.Ida); k++ {
								if k == len(reserva.Passengers[i].Ancillaries.Ida) {
									seleccionAncillary.Cantidad = 1
									reserva.Passengers[i].Ancillaries.Ida = append(reserva.Passengers[i].Ancillaries.Ida, seleccionAncillary)
									break
								}

								if seleccionAncillary.SSR == reserva.Passengers[i].Ancillaries.Ida[k].SSR {
									reserva.Passengers[i].Ancillaries.Ida[k].Cantidad += 1
									break
								}

							}

						}

						if len(reserva.Vuelos) == 2 {

							fmt.Println("\nAncillares vuelta:")

							fmt.Print("\nIngrese los Ancillaries (separados por comas): ")

							fmt.Scan(&seleccionAncillaries)

							stringSlice = strings.Split(seleccionAncillaries, ",")
							seleccionArray = make([]int, len(stringSlice))

							for j, v := range stringSlice {
								seleccionArray[j], err = strconv.Atoi(v)
								seleccionArray[j] -= 1

								var seleccionAncillary models.PassengerAncillary

								seleccionAncillary.SSR = vuelos[0].Ancillaries[seleccionArray[j]].SSR

								for k := 0; k <= len(reserva.Passengers[i].Ancillaries.Vuelta); k++ {
									if k == len(reserva.Passengers[i].Ancillaries.Vuelta) {
										seleccionAncillary.Cantidad = 1
										reserva.Passengers[i].Ancillaries.Vuelta = append(reserva.Passengers[i].Ancillaries.Vuelta, seleccionAncillary)
										break
									}

									if seleccionAncillary.SSR == reserva.Passengers[i].Ancillaries.Vuelta[k].SSR {
										reserva.Passengers[i].Ancillaries.Vuelta[k].Cantidad += 1
										break
									}

								}

							}

						}

					}

					JSONString, err := json.Marshal(reserva)
					if err != nil {
						fmt.Println(err)
						return
					}

					req, err := http.NewRequest("PUT", URL+"/reserva?pnr="+reserva.PNR+"&apellido="+reserva.Apellido, bytes.NewBuffer(JSONString))

					if err != nil {
						// Handle error
					}

					client := &http.Client{}
					resp, err = client.Do(req)

					if err != nil {
						// Handle error
					}
					defer resp.Body.Close()

				case 3:
					return

				default:
					fmt.Println("Opcion invalida")
				}
			case 4:
				fmt.Println("Salir")
				return
			default:
				fmt.Println("Opcion invalida")
			}
		case 2:
			fmt.Println("Obtener estadisticas")
			// Do something
		case 3:
			fmt.Println("Salir")
			return
		default:
			fmt.Println("Opcion invalida")
		}
	}
}
