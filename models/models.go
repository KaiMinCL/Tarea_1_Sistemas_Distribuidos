// This package defines several structs that represent different entities in an airline reservation system,
// and provides functions to interact with a MongoDB database to perform CRUD operations on flight reservations.

package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//STRUCT DEFINITIONS

// Define the Plane struct, which represents an airplane model
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
	Name        string                   `bson:"name" json:"name"`
	Apellido    string                   `bson:"apellido" json:"apellido"`
	Edad        int                      `bson:"edad" json:"edad"`
	Ancillaries []PassengerAncillaryList `bson:"ancillaries" json:"ancillaries"`
}

// Define the Reservations struct, which represents a reservation made by one or more passengers for one or more flights
type Reservation struct {
	PNR        string      `bson:"PNR" json:"PNR"`
	Apellido   string      `bson:"apellido" json:"apellido"`
	Vuelos     []Flight    `bson:"vuelos" json:"vuelos"`
	Passengers []Passenger `bson:"pasajeros" json:"pasajeros"`
}

// Connect to MongoDB and retrieve the collection needed
func getDatabaseCollection(collectionName string) *mongo.Collection {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var CONNECTION_STRING = os.Getenv("CONNECTION_STRING")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(CONNECTION_STRING))

	if err != nil {
		panic(err)
	}

	collection := client.Database("aerolinea").Collection(collectionName)

	return collection
}

//CRUD Vuelos

// This function retrieves one document from the "vuelos" collection of the "aerolinea" database and stores it in a Flight struct

func GetVuelos(origenVuelo string, destinoVuelo string, fechaVuelo string) ([]Flight, error) {

	var vuelos []Flight

	collection := getDatabaseCollection("vuelos")

	// Define a filter to find all flights with the specified origen, destino, and fecha
	filter := bson.M{"origen": origenVuelo, "destino": destinoVuelo, "fecha": fechaVuelo}

	// Find all the flights that match the filter
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Error retrieving documents:", err)
		return vuelos, err
	}
	defer cur.Close(context.Background())

	// Iterate over the cursor and append each flight to the vuelos slice
	for cur.Next(context.Background()) {
		var vuelo Flight
		if err := cur.Decode(&vuelo); err != nil {
			fmt.Println("Error decoding document:", err)
			return vuelos, err
		}
		vuelos = append(vuelos, vuelo)
	}

	return vuelos, nil
}

// This function deletes one document from the "vuelos" collection of the "aerolinea" database=
func DeleteVuelo(numeroVuelo string, origenVuelo string, destinoVuelo string, fechaVuelo string) error {

	collection := getDatabaseCollection("vuelos")
	result, err := collection.DeleteOne(context.Background(), bson.M{"numero_vuelo": numeroVuelo, "origen": origenVuelo, "destino": destinoVuelo, "fecha": fechaVuelo})
	if err != nil {
		// Handle the error
		fmt.Println("Error deleting document:", err)
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("flight not found")
	}

	return nil
}

// This function inserts one document from the "vuelos" collection of the "aerolinea" database
func CreateVuelo(flight Flight) {
	collection := getDatabaseCollection("vuelos")
	_, err := collection.InsertOne(context.Background(), flight)
	if err != nil {
		// Handle the error
		fmt.Println("Error inserting document:", err)
	}
	return
}

// This function updates one document from the "vuelos" collection of the "aerolinea" database
func UpdateVuelo(numeroVuelo string, origenVuelo string, destinoVuelo string, fechaVuelo string, stockDePasajeros int) (map[string]interface{}, error) {

	collection := getDatabaseCollection("vuelos")

	filter := bson.M{"numero_vuelo": numeroVuelo, "origen": origenVuelo, "destino": destinoVuelo, "fecha": fechaVuelo}

	updateBSON := bson.M{"$set": bson.M{"avion.stock_de_pasajeros": stockDePasajeros}}
	_, err := collection.UpdateOne(context.Background(), filter, updateBSON)
	if err != nil {
		// Handle the error
		fmt.Println("Error updating document:", err)
		return nil, err
	}

	fmt.Println("Document updated successfully")

	updatedDoc := Flight{}
	err = collection.FindOne(context.Background(), filter).Decode(&updatedDoc)

	if err != nil {
		// Handle the error
		fmt.Println("Error decoding updated document:", err)
		return nil, err
	}

	response := map[string]interface{}{
		"numero_vuelo": updatedDoc.NumeroVuelo,
		"origen":       updatedDoc.Origen,
		"destino":      updatedDoc.Destino,
		"hora_salida":  updatedDoc.HoraSalida,
		"hora_llegada": updatedDoc.HoraLlegada,
	}

	return response, nil
}

// CRUD Reservations

// CreateReservation adds a new reservation to the database
func CreateReservation(reservation Reservation) error {
	collection := getDatabaseCollection("reservas")
	_, err := collection.InsertOne(context.Background(), reservation)
	return err
}

// GetReservation returns a reservation from the database with the given ID
func GetReservation(id primitive.ObjectID) (Reservation, error) {
	collection := getDatabaseCollection("reservas")
	var reservation Reservation
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&reservation)
	return reservation, err
}

// UpdateReservation updates a reservation in the database with the given ID
func UpdateReservation(id primitive.ObjectID, reservation Reservation) error {
	collection := getDatabaseCollection("reservas")
	_, err := collection.ReplaceOne(context.Background(), bson.M{"_id": id}, reservation)
	return err
}

// DeleteReservation deletes a reservation from the database with the given ID
func DeleteReservation(id primitive.ObjectID) error {
	collection := getDatabaseCollection("reservas")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
