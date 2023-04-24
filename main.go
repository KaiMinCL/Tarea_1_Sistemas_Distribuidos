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
		clearScreen()

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
				// - check if there are seats left on the plane before selling them
				// - Remove capacity for every passenger sold on the flights
				// - Remove capacity for every ancillary sold on the flight
				// - Update mongoDB with the updated flight information

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

					// Check if the vuelos slice is empty
					if len(vuelos) == 0 {
						fmt.Println("No hay vuelos disponibles para la fecha de vuelta")
						break
					}

					fmt.Println("Vuelta:")

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

						precioVuelo := 590 * minutosVuelo
						fmt.Print("\t", i+1)
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
						fmt.Print(". " + vueloIda.Ancillaries[i].Nombre + " $" + fmt.Sprint(AncillaryPrice[vueloIda.Ancillaries[i].SSR]) + "\n")
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
							time.Sleep(50 * time.Millisecond)

							fmt.Print("\t", i+1)
							fmt.Print(". " + vueloVuelta.Ancillaries[i].Nombre + " $" + fmt.Sprint(AncillaryPrice[vueloIda.Ancillaries[i].SSR]) + "\n")
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

				waitAnimation()
				clearScreen()

				fmt.Println("\nLa reserva fue generada con el PNR:", Reserva.PNR)

				var costoTotal int

				for _, pasajero := range Reserva.Passengers {
					costoTotal += pasajero.Balances.AncillariesIda + pasajero.Balances.VueloIda + pasajero.Balances.VueloVuelta + pasajero.Balances.VueloVuelta
				}

				fmt.Println("\nEl costo total de la reserva fue de: $" + fmt.Sprint(costoTotal))

			case 2:

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

					waitAnimation()
					clearScreen()

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

					// Check if the vuelos slice is empty
					if len(vuelos) == 0 {
						fmt.Println("No hay vuelos disponibles para la fecha indicada")
						break
					}

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

						precioVuelo := 590 * minutosVuelo
						fmt.Print(i + 1)
						fmt.Print(". " + vuelos[i].NumeroVuelo + " " + vuelos[i].HoraSalida + " - " + vuelos[i].HoraLlegada + " $")
						fmt.Print(precioVuelo, "\n")
					}

					var flightOption int

					fmt.Print("\nIngrese una Opción: ")
					fmt.Scan(&flightOption)
					time.Sleep(25 * time.Millisecond)

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

					waitAnimation()
					clearScreen()

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

					fmt.Println("¡La reserva fue modificada exitosamente!")

					fmt.Print("\nPresione cualquier tecla para volver al menu principal...")
					bufio.NewReader(os.Stdin).ReadString('\n')

					var input string
					fmt.Scanln(&input)

				case 2:

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

					fmt.Println("Ingrese una opción: ")
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
						fmt.Print(i + 1)
						fmt.Print(". " + vuelos[0].Ancillaries[i].Nombre + " $" + fmt.Sprint(AncillaryPrice[vuelos[0].Ancillaries[i].SSR]) + "\n")
					}

					fmt.Println("Pasajeros: ")
					for i := 0; i < len(reserva.Passengers); i++ {
						time.Sleep(50 * time.Millisecond)

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
			for mes, valor := range data["promedio_pasajeros"].(map[string]interface{}) {
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
