// This package defines several structs that represent different entities in an airline reservation system,
// and provides functions to interact with a MongoDB database to perform CRUD operations on flight reservations.

package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define the Plane struct, which represents an airplane model
type Plane struct {
	Modelo           string `bson:"modelo"`
	NumeroDeSerie    string `bson:"numero_de_serie"`
	StockDePasajeros int    `bson:"stock_de_pasajeros"`
}

// Define the FlightAncillary struct, which represents an optional service that can be offered on a flight
type FlightAncillary struct {
	Nombre string `bson:"nombre"`
	Stock  int    `bson:"stock"`
	SSR    string `bson:"ssr"`
}

// Define the Flight struct, which represents a flight that can be booked by passengers
type Flight struct {
	NumeroVuelo string            `bson:"numero_vuelo"`
	Origen      string            `bson:"origen"`
	Destino     string            `bson:"destino"`
	HoraSalida  string            `bson:"hora_salida"`
	HoraLlegada string            `bson:"hora_llegada"`
	Fecha       string            `bson:"fecha"`
	Avion       Plane             `bson:"avion"`
	Ancillaries []FlightAncillary `bson:"ancillaries"`
}

// Define the PassengerAncillary struct, which represents an optional service that a passenger can select
type PassengerAncillary struct {
	Cantidad int    `bson:"cantidad"`
	SSR      string `bson:"ssr"`
}

// Define the PassengerAncillaryList struct, which represents the ancillaries selected by a passenger for a round-trip flight
type PassengerAncillaryList struct {
	ida    []PassengerAncillary `bson:"ida"`
	vuelta []PassengerAncillary `bson:"vuelta"`
}

// Define the Passengers struct, which represents a passenger who has booked a flight
type Passengers struct {
	Name        string                   `bson:"name"`
	Apellido    string                   `bson:"apellido"`
	Edad        int                      `bson:"edad"`
	Ancillaries []PassengerAncillaryList `bson:"ancillaries"`
}

// Define the Reservations struct, which represents a reservation made by one or more passengers for one or more flights
type Reservations struct {
	Vuelos     []Flight     `bson:"vuelos"`
	Passengers []Passengers `bson:"pasajeros"`
}

// This function takes a JSON string as input and returns a pretty formatted JSON string with proper indentation
func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
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

// This function retrieves all documents in the "vuelos" collection of the "aerolinea" database and prints them in pretty formatted JSON
func GetVuelos() {

	collection := getDatabaseCollection("vuelos")

	// Find all documents in the "vuelos" collection
	cursor, err := collection.Find(context.Background(), bson.D{})

	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var print = ""
	// Iterate over the cursor and decode each document
	for cursor.Next(context.Background()) {
		var vuelo Flight
		if err := cursor.Decode(&vuelo); err != nil {
			log.Fatal(err)
		}

		b, err := json.Marshal(vuelo)
		if err != nil {
			fmt.Println(err)
			return
		}

		stringVuelo, _ := PrettyString(string(b))
		print += stringVuelo
	}

	fmt.Println("[\n" + print + "\n]")
}

// This function retrieves one document from the "vuelos" collection of the "aerolinea" database and stores it in a Flight struct

func GetVuelo(id primitive.ObjectID) {

	var vuelo Flight

	collection := getDatabaseCollection("vuelos")
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&vuelo)
	if err != nil {
		// Handle the error
		fmt.Println("Error retrieving document:", err)
		return
	}

}

// This function deletes one document from the "vuelos" collection of the "aerolinea" database

func DeleteVuelo(id primitive.ObjectID) {

	collection := getDatabaseCollection("vuelos")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		// Handle the error
		fmt.Println("Error deleting document:", err)
		return
	}

}

// This function inserts one document from the "vuelos" collection of the "aerolinea" database
func CreateVuelo(flight Flight) {
	collection := getDatabaseCollection("vuelos")
	_, err := collection.InsertOne(context.Background(), flight)
	if err != nil {
		// Handle the error
		fmt.Println("Error inserting document:", err)
		return
	}

}

// This function updates one document from the "vuelos" collection of the "aerolinea" database
func UpdateVuelo(id primitive.ObjectID) {

	collection := getDatabaseCollection("vuelos")

	filter := bson.D{}
	update := bson.D{}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		// Handle the error
		fmt.Println("Error updating document:", err)
		return
	}

}

// CreateReservation adds a new reservation to the database
func CreateReservation(reservation Reservations) error {
	collection := getDatabaseCollection("reservas")
	_, err := collection.InsertOne(context.Background(), reservation)
	return err
}

// GetReservation returns a reservation from the database with the given ID
func GetReservation(id primitive.ObjectID) (Reservations, error) {
	collection := getDatabaseCollection("reservas")
	var reservation Reservations
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&reservation)
	return reservation, err
}

// UpdateReservation updates a reservation in the database with the given ID
func UpdateReservation(id primitive.ObjectID, reservation Reservations) error {
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
