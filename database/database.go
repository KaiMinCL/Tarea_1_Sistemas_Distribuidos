package database

import (
	"bd_aerolinea/models"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func GetVuelos(origenVuelo string, destinoVuelo string, fechaVuelo string) ([]models.Flight, error) {

	var vuelos []models.Flight

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
		var vuelo models.Flight
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
func CreateVuelo(flight models.Flight) {
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

	updatedDoc := models.Flight{}
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

func CreateReservation(reservation models.Reservation) (map[string]interface{}, error) {

	collection := getDatabaseCollection("reservas")
	_, err := collection.InsertOne(context.Background(), reservation)
	response := map[string]interface{}{
		"PNR": reservation.PNR,
	}

	return response, err
}

// GetReservation returns a reservation from the database with the given ID
func GetReservation(pnr string, apellido string) (models.Reservation, error) {
	collection := getDatabaseCollection("reservas")
	var reservation models.Reservation
	err := collection.FindOne(context.Background(), bson.M{"PNR": pnr, "apellido": bson.M{"$regex": apellido, "$options": "i"}}).Decode(&reservation)
	return reservation, err
}

func GetAllReservations() ([]models.Reservation, error) {
	var reservas []models.Reservation

	collection := getDatabaseCollection("reservas")
	filter := bson.M{}
	// Find all the reservations
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Error retrieving documents:", err)
		return reservas, err
	}
	defer cur.Close(context.Background())

	// Iterate over the cursor and append each reservation to the array
	for cur.Next(context.Background()) {
		var reserva models.Reservation
		if err := cur.Decode(&reserva); err != nil {
			fmt.Println("Error decoding document:", err)
			return reservas, err
		}
		reservas = append(reservas, reserva)
	}

	return reservas, nil
}

// UpdateReservation updates a reservation in the database with the given ID
func UpdateReservation(pnr string, apellido string, reservation models.Reservation) (map[string]interface{}, error) {
	collection := getDatabaseCollection("reservas")
	_, err := collection.ReplaceOne(context.Background(), bson.M{"PNR": pnr, "apellido": apellido}, reservation)

	response := map[string]interface{}{
		"PNR": reservation.PNR,
	}

	return response, err
}

// DeleteReservation deletes a reservation from the database with the given ID
func DeleteReservation(pnr string, apellido string) error {
	collection := getDatabaseCollection("reservas")
	_, err := collection.DeleteOne(context.Background(), bson.M{"PNR": pnr, "apellido": apellido})
	return err
}
