package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"github.com/joho/godotenv"
	"encoding/json"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		SERVER = os.Getenv("SERVER")
		PORT = os.Getenv("PORT")

		option int
		subOption int
		fechaIda string
		fechaRegreso string
		origen string
		destino string
		cantidadPasajeros int
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
				fmt.Println("Crear reserva")
				fmt.Println("Ingrese a fecha de ida:")
				fmt.Scan(&fechaIda)
				fmt.Println("Ingrese a fecha de regreso:")
				fmt.Scan(&fechaRegreso)
				fmt.Println("Ingrese origen:")
				fmt.Scan(&origen)
				fmt.Println("Ingrese destino")
				fmt.Scan(&destino)
				fmt.Println("Cantidad de Pasajeros:")
				fmt.Scan(&cantidadPasajeros)
				fmt.Println("Vuelos disponibles:")
				fmt.Println("Ida:")

				resp, err := http.Get("/api/vuelos?origen=" + origen +"&destino=" + destino + "&fecha=" + fechaIda)

				if err != nil{
					log.Fatal(err)
				}

				var vuelos []models.Flight

				defer resp.Body.Close()

				json.NewDecoder(resp.Body).Decode(vuelos)

				var count int = 1
				for i := 0; i<len(vuelos); i++{
					if vuelos[i].Avion.StockDePasajeros <= cantidadPasajeros{
						fmt.Printf("%v. %v %v - %v", count, vuelos[i].NumeroVuelo, vuelos[i].HoraSalida, vuelos[i].HoraLlegada)
						count++
					}
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
