// This package defines several structs that represent different entities in an airline reservation system,
// and provides functions to interact with a MongoDB database to perform CRUD operations on flight reservations.

package models

import (
	"fmt"
	"time"
)

//STRUCT DEFINITIONS

// Define the Plane struct, which represents an airplane mode/

type Plane struct {
	Modelo           string `bson:"modelo" json:"modelo"`
	NumeroDeSerie    string `bson:"numero_de_serie" json:"numero_de_serie"`
	StockDePasajeros int    `bson:"stock_de_pasajeros" json:"stock_de_pasajeros"`
}

// Define the FlightAncillary struct, which represents an optional service that can be offered on a flight
type FlightAncillary struct {
	Nombre string `bson:"nombre" json:"nombre"`
	Stock  int    `bson:"stock" json:"stock"`
	SSR    string `bson:"ssr" json:"ssr"`
}

// Define the Flight struct, which represents a flight that can be booked by passengers
type Flight struct {
	NumeroVuelo string            `bson:"numero_vuelo" json:"numero_vuelo"`
	Origen      string            `bson:"origen" json:"origen"`
	Destino     string            `bson:"destino" json:"destino"`
	HoraSalida  string            `bson:"hora_salida" json:"hora_salida"`
	HoraLlegada string            `bson:"hora_llegada" json:"hora_llegada"`
	Fecha       string            `bson:"fecha" json:"fecha"`
	Avion       Plane             `bson:"avion" json:"avion"`
	Ancillaries []FlightAncillary `bson:"ancillaries" json:"ancillaries"`
}

type ReservationFlight struct {
	NumeroVuelo string `bson:"numero_vuelo" json:"numero_vuelo"`
	Origen      string `bson:"origen" json:"origen"`
	Destino     string `bson:"destino" json:"destino"`
	HoraSalida  string `bson:"hora_salida" json:"hora_salida"`
	HoraLlegada string `bson:"hora_llegada" json:"hora_llegada"`
	Fecha       string `bson:"fecha" json:"fecha"`
}

// Define the PassengerAncillary struct, which represents an optional service that a passenger can select
type PassengerAncillary struct {
	Cantidad int    `bson:"cantidad" json:"cantidad"`
	SSR      string `bson:"ssr" json:"ssr"`
}

// Define the PassengerAncillaryList struct, which represents the ancillaries selected by a passenger for a round-trip flight
type PassengerAncillaryList struct {
	Ida    []PassengerAncillary `bson:"ida" json:"ida"`
	Vuelta []PassengerAncillary `bson:"vuelta" json:"vuelta"`
}

// Define the Passengers struct, which represents a passenger who has booked a flight
type Passenger struct {
	Name        string                 `bson:"nombre" json:"nombre"`
	Apellido    string                 `bson:"apellido" json:"apellido"`
	Edad        int                    `bson:"edad" json:"edad"`
	Ancillaries PassengerAncillaryList `bson:"ancillaries" json:"ancillaries"`
	Balances    BalancesObj            `bson:"balances" json:"balances"`
}

type BalancesObj struct {
	AncillariesIda    int `bson:"ancillaries_ida" json:"ancillaries_ida"`
	VueloIda          int `bson:"vuelo_ida" json:"vuelo_ida"`
	AncillariesVuelta int `bson:"ancillaries_vuelta" json:"ancillaries_vuelta"`
	VueloVuelta       int `bson:"vuelo_vuelta" json:"vuelo_vuelta"`
}

// Define the Reservations struct, which represents a reservation made by one or more passengers for one or more flights

type Reservation struct {
	PNR        string              `bson:"PNR" json:"PNR"`
	Apellido   string              `bson:"apellido" json:"apellido"`
	Vuelos     []ReservationFlight `bson:"vuelos" json:"vuelos"`
	Passengers []Passenger         `bson:"pasajeros" json:"pasajeros"`
}

type AncillariesStatistics struct {
	Nombre   string `bson:"nombre" json:"nombre"`
	SSR      string `bson:"ssr" json:"ssr"`
	Ganancia int    `bson:"ganancia" json:"ganancia"`
}

type PassengerAverage struct {
	Jan int `bson:"enero" json:"enero"`
	Feb int `bson:"febrero" json:"febrero"`
	Mar int `bson:"marzo" json:"marzo"`
	Apr int `bson:"abril" json:"abril"`
	May int `bson:"mayo" json:"mayo"`
	Jun int `bson:"junio" json:"junio"`
	Jul int `bson:"julio" json:"julio"`
	Aug int `bson:"agosto" json:"agosto"`
	Sep int `bson:"septiembre" json:"septiembre"`
	Oct int `bson:"octubre" json:"octubre"`
	Nov int `bson:"noviembre" json:"noviembre"`
	Dec int `bson:"diciembre" json:"diciembre"`
}

type Statistics struct {
	RutaMayorGanancia  string                  `bson:"ruta_mayor_ganancia" json:"ruta_mayor_ganancia"`
	RutaMenorGanancia  string                  `bson:"ruta_menor_ganancia" json:"ruta_menor_ganancia"`
	RankingAncillaries []AncillariesStatistics `bson:"ranking_ancillaries" json:"ranking_ancillaries"`
	PromedioPasajeros  PassengerAverage        `bson:"promedio_pasajeros" json:"promedio_pasajeros"`
}

type PNRCapsule struct {
	PNR string `bson:"PNR" json:"PNR"`
}

var (
	AncillaryPrice = map[string]int{"BGH": 10000, "BGR": 30000, "STDF": 5000, "PAXS": 2000, "PTCR": 40000, "AVIH": 40000, "SPML": 35000, "LNGE": 15000, "WIFI": 20000}
	AncillaryName  = map[string]string{"BGH": "Equipaje de mano", "BGR": "Equipaje de bodega", "STDF": "Asiento", "PAXS": "Embarque y Check In prioritario", "PTCR": "Mascota en cabina", "AVIH": "Mascota en bodega", "SPML": "Equipaje especial", "LNGE": "Acceso a Salon VIP", "WIFI": "Wi-Fi a bordo"}
)

func SumAncillaries(passangerAncillaries PassengerAncillaryList) (int, int) {
	//This function calculates the price for the ancillaries associated with this passenger.
	countIda := 0
	countVuelta := 0
	for j := 0; j < len(passangerAncillaries.Ida); j++ {
		countIda += AncillaryPrice[passangerAncillaries.Ida[j].SSR] * passangerAncillaries.Ida[j].Cantidad
	}
	for j := 0; j < len(passangerAncillaries.Vuelta); j++ {
		countVuelta += AncillaryPrice[passangerAncillaries.Vuelta[j].SSR] * passangerAncillaries.Vuelta[j].Cantidad
		//fm.Println(countVuelta, PA.Vuelta[j].SSR, PA.Vuelta[j].Cantidad)t
	}

	return countIda, countVuelta
}

func SumVuelos(flights []ReservationFlight) (int, int) {
	//This function returns the sum of the price of all the tickets associated with this passenger.
	var tiempoIda int = 0
	var tiempoVuelta int = 0

	horaSalidaIda, _ := time.Parse("15:04", flights[0].HoraSalida)
	horaLlegadaIda, _ := time.Parse("15:04", flights[0].HoraLlegada)

	if horaLlegadaIda.Before(horaSalidaIda) {
		horaLlegadaIda = horaLlegadaIda.Add(24 * time.Hour)
	}

	tiempoIda = int(horaLlegadaIda.Sub(horaSalidaIda).Minutes())

	if len(flights) == 2 {
		horaSalidaVuelta, _ := time.Parse("15:04", flights[1].HoraSalida)
		horaLlegadaVuelta, _ := time.Parse("15:04", flights[1].HoraLlegada)

		if horaLlegadaVuelta.Before(horaSalidaVuelta) {
			horaLlegadaVuelta = horaLlegadaVuelta.Add(24 * time.Hour)
		}

		tiempoVuelta = int(horaLlegadaVuelta.Sub(horaSalidaVuelta).Minutes())
	} else if len(flights) > 2 {
		fmt.Println("No deben haber más de dos vuelos por reserva")
	}

	return tiempoIda * 590, tiempoVuelta * 590
}
