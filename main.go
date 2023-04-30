package main

import (
	"bd_aerolinea/flightbooking"
	"bd_aerolinea/models"
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
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

	for true {
		flightbooking.ClearScreen()
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
			flightbooking.ClearScreen()

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
				flightbooking.ClearScreen()
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
					pasajeros         []models.Passenger
					reserva           models.Reservation
					vueloIda          models.Flight
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

				flightbooking.WaitAnimation()
				flightbooking.ClearScreen()

				vuelos, err := flightbooking.CheckFlightAvailability(URL, origen, destino, fechaIda)
				if err != nil {
					log.Fatal(err)
					break
				}

				// Revisar si hay vuelos disponibles
				if len(vuelos) == 0 {
					fmt.Println("No hay vuelos disponibles para la fecha indicada")
					break
				}

				flightbooking.DisplayFlights(vuelos, "Ida: ", false)

				vueloIda = flightbooking.SelectFlight(vuelos, &reserva, origen, destino, fechaIda, false)

				vueloIda.Avion.StockDePasajeros -= 1

				if fechaRegreso != "no" {

					vuelos, err = flightbooking.CheckFlightAvailability(URL, destino, origen, fechaRegreso)
					if err != nil {
						log.Fatal(err)
						break
					}

					flightbooking.DisplayFlights(vuelos, "Vuelta: ", false)

					vueloVuelta = flightbooking.SelectFlight(vuelos, &reserva, destino, origen, fechaRegreso, false)

					vueloVuelta.Avion.StockDePasajeros -= 1

					vuelos = []models.Flight{vueloIda, vueloVuelta}

					pasajeros, err = flightbooking.PassengerInformation(cantidadPasajeros, vuelos)
				} else {
					vuelos = []models.Flight{vueloIda}
					pasajeros, err = flightbooking.PassengerInformation(cantidadPasajeros, vuelos)
				}

				reserva.Passengers = pasajeros
				reserva.Apellido = pasajeros[0].Apellido

				PNR, err := flightbooking.GeneratePNR(URL)
				if err != nil {
					log.Fatal(err)
				}
				reserva.PNR = PNR

				err = flightbooking.MakeReservation(URL, reserva, false)
				if err != nil {
					log.Fatal(err)
				}

				err = flightbooking.UpdateFlightStock(URL, vueloIda)
				if err != nil {
					fmt.Println(err)
					return
				}

				if fechaRegreso != "no" {
					err = flightbooking.UpdateFlightStock(URL, vueloVuelta)
					if err != nil {
						fmt.Println(err)
						return
					}
				}
				flightbooking.WaitAnimation()
				flightbooking.ClearScreen()

				fmt.Println("\nLa reserva fue generada con el PNR:", reserva.PNR)

				var costoTotal int

				for _, pasajero := range reserva.Passengers {
					costoTotal += pasajero.Balances.AncillariesIda + pasajero.Balances.VueloIda + pasajero.Balances.VueloVuelta + pasajero.Balances.AncillariesVuelta
				}

				fmt.Println("\nEl costo total de la reserva fue de: $" + fmt.Sprint(costoTotal))

				fmt.Print("\nPresione cualquier tecla para volver al menu principal...")
				bufio.NewReader(os.Stdin).ReadString('\n')

				var input string
				fmt.Scanln(&input)

			case 2:

				// ALL OF THIS IS FOR VIEWING A RESERVATION

				var reserva models.Reservation

				fmt.Print("Ingrese el PNR: ")
				fmt.Scan(&reserva.PNR)
				fmt.Print("Ingrese el Apellido: ")
				fmt.Scan(&reserva.Apellido)

				err = flightbooking.GetReserva(URL, &reserva)
				if err != nil {
					log.Fatal(err)
				}

				flightbooking.WaitAnimation()
				flightbooking.ClearScreen()

				fmt.Println("Reserva: ")

				fmt.Println("Vuelo Ida: ")

				fmt.Println("\tNro: " + reserva.Vuelos[0].NumeroVuelo + " | Horario: " + reserva.Vuelos[0].HoraSalida + " - " + reserva.Vuelos[0].HoraLlegada)

				if len(reserva.Vuelos) == 2 {
					fmt.Println("Vuelo vuelta: ")
					fmt.Println("\t Nro: " + reserva.Vuelos[1].NumeroVuelo + " | Horario: " + reserva.Vuelos[1].HoraSalida + " - " + reserva.Vuelos[1].HoraLlegada)
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

				var reserva models.Reservation

				fmt.Print("Ingrese el PNR: ")
				fmt.Scan(&reserva.PNR)

				fmt.Print("Ingrese el Apellido: ")
				fmt.Scan(&reserva.Apellido)

				err = flightbooking.GetReserva(URL, &reserva)
				if err != nil {
					fmt.Println("No existe la reserva")
					break
				}

				flightbooking.WaitAnimation()
				flightbooking.ClearScreen()

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

					flightbooking.WaitAnimation()
					flightbooking.ClearScreen()

					var (
						flightSelected int
						fecha          string
					)

					fmt.Println("Vuelos: ")
					time.Sleep(50 * time.Millisecond)

					fmt.Println("\t1. Ida: " + reserva.Vuelos[0].NumeroVuelo + " " + reserva.Vuelos[0].HoraSalida + " - " + reserva.Vuelos[0].HoraLlegada)

					if len(reserva.Vuelos) == 2 {
						time.Sleep(50 * time.Millisecond)
						fmt.Println("\t2. Vuelta: " + reserva.Vuelos[1].NumeroVuelo + " " + reserva.Vuelos[1].HoraSalida + " - " + reserva.Vuelos[1].HoraLlegada)
					}

					fmt.Print("Ingrese una opción: ")
					fmt.Scan(&flightSelected)
					time.Sleep(25 * time.Millisecond)

					flightSelected -= 1

					fmt.Print("Ingrese nueva fecha: ")
					fmt.Scan(&fecha)
					time.Sleep(25 * time.Millisecond)

					vuelos, err := flightbooking.CheckFlightAvailability(URL, reserva.Vuelos[flightSelected].Origen, reserva.Vuelos[flightSelected].Destino, fecha)
					if err != nil {
						log.Fatal(err)
						break
					}

					// Revisar si hay vuelos disponibles
					if len(vuelos) == 0 {
						fmt.Println("No hay vuelos disponibles para la fecha indicada")
						break
					}

					flightbooking.WaitAnimation()
					flightbooking.ClearScreen()

					vueloViejo := reserva.Vuelos[flightSelected]

					// Si se encuentra el vuelo en la fecha, eliminarlo de las opciones
					for i, vuelo := range vuelos {
						if vuelo.NumeroVuelo == vueloViejo.NumeroVuelo {
							vuelos = append(vuelos[:i], vuelos[i+1:]...)
						}
					}

					flightbooking.DisplayFlights(vuelos, "", true)

					// Se elimina el vuelo viejo de mis reservas
					for i, _ := range reserva.Vuelos {

						if vueloViejo.NumeroVuelo == reserva.Vuelos[flightSelected].NumeroVuelo {
							reserva.Vuelos = append(reserva.Vuelos[:i], reserva.Vuelos[i+1:]...)
						}

					}

					nuevoVuelo := flightbooking.SelectFlight(vuelos, &reserva, vueloViejo.Origen, vueloViejo.Destino, fecha, true)

					nuevoVuelo.Avion.StockDePasajeros -= len(reserva.Passengers)

					err = flightbooking.MakeReservation(URL, reserva, true)
					if err != nil {
						fmt.Println("Error:", err)
						break
					}

					vuelos, err = flightbooking.GetVuelos(URL, vueloViejo.Origen, vueloViejo.Destino, vueloViejo.Fecha)
					if err != nil {
						fmt.Println("Error:", err)
						break
					}

					// Se obtiene el vuelo viejo para aumentar nuevamente el stock de pasajeros
					var updateVueloViejo models.Flight
					for _, vuelo := range vuelos {
						if vuelo.NumeroVuelo == vueloViejo.NumeroVuelo {
							updateVueloViejo = vuelo
						}
					}

					updateVueloViejo.Avion.StockDePasajeros += len(reserva.Passengers)

					err = flightbooking.UpdateFlightStock(URL, updateVueloViejo)
					if err != nil {
						fmt.Println("Error:", err)
						break
					}
					err = flightbooking.UpdateFlightStock(URL, nuevoVuelo)
					if err != nil {
						fmt.Println("Error:", err)
						break
					}

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

					flightbooking.WaitAnimation()
					flightbooking.ClearScreen()

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

					vuelos, err := flightbooking.GetVuelos(URL, reserva.Vuelos[flightOption].Origen, reserva.Vuelos[flightOption].Destino, reserva.Vuelos[flightOption].Fecha)
					if err != nil {
						fmt.Println("error:", err)
						break
					}

					// Selecciona solo el vuelo que buscamos
					for _, vuelo := range vuelos {
						if vuelo.NumeroVuelo == reserva.Vuelos[flightOption].NumeroVuelo {
							vuelos = []models.Flight{vuelo}
						}
					}

					flightbooking.WaitAnimation()

					fmt.Println("\nPasajeros: ")

					for i := 0; i < len(reserva.Passengers); i++ {

						fmt.Println("\t", reserva.Passengers[i].Name, reserva.Passengers[i].Edad)

						selectedAncillaries, err := flightbooking.SelectAncillaries("", vuelos[0].Ancillaries)
						if err != nil {
							fmt.Println(err)
							break
						}

						flightbooking.AddAncillaries(selectedAncillaries, flightOption, i, &reserva)
					}

					err = flightbooking.MakeReservation(URL, reserva, true)
					if err != nil {
						log.Fatal(err)
						fmt.Println("Ha ocurrido un error realizando la reserva")
					}

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
			flightbooking.DisplayStatistics(URL)

		case 3:
			flightbooking.ClearScreen()
			return
		default:
			fmt.Println("Opcion invalida")
		}
	}
	flightbooking.ClearScreen()
}
