package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var SERVER = os.Getenv("SERVER")
	var PORT = os.Getenv("PORT")

	fmt.Println("connecting to", SERVER, ":", PORT)
	var option int
	var subOption int
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
				// Do something
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
