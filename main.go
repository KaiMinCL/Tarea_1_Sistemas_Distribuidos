package main

import (
	"bd_aerolinea/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
				var (
					fechaIda          string
					fechaRegreso      string
					origen            string
					destino           string
					cantidadPasajeros int
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
					fmt.Println("There are no vuelos available.")
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
					fmt.Print(i)
					fmt.Print(". " + vuelos[i].NumeroVuelo + " " + vuelos[i].HoraSalida + " - " + vuelos[i].HoraLlegada + " $")
					fmt.Print(precioVuelo, "\n")
				}

			case 2:
				fmt.Println("Obtener reserva")
				// Do something
			case 3:
				fmt.Println("Modificar reserva")
				// Do something
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
