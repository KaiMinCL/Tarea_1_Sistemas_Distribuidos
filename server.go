package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"context"
	"log"
	"fmt"
)

const(
	DB_NAME = "aerolinea"
)

type plane struct {
	modelo string `json:"modelo"`
	numero_de_serie string `json:"numero_de_serie"`
	stock_de_pasajeros int `json:"stock_de_pasajeros"`
}

type ancillary struct {
	nombre string `json:"nombre"`
	stock int `json:"stock"`
	ssr string `json:"ssr"`
}

type passenger_ancillary struct{
	DepartureFlight []ancillary `json:"departureFlight"`
	ReturnFlight []ancillary `json:"returnFlight"`
}

type flight struct{
	id primitive.ObjectID `json:"_id, omitempty"`
	numero_vuelo string `json:"numero_vuelo"`
	origen string `json:"origen"`
	destino string `json:"destino"`
	hora_salida string `json:"hora_salida"`
	hora_llegada string `json:"hora_llegada"`
	fecha string `json:"fecha"`
	avion plane `json:"avion"`
	ancillaries []ancillary `json:"ancillaries"`
}

type balance struct{
	AncillariesDepartureFlight int `json:"ancillariesDepartureFlight"`
	DepartureFlight int `json:"departureFlight"`
	AncillariesReturnFlight int `json:"ancillariesReturnFlight"`
	ReturnFlight int `json:"returnFlight"`
}

type passenger struct{
	FirstName string `json:"firstName"`
	Surname string `json:"surname"`
	Age int `json:"age"`
	Ancillaries passenger_ancillary `json:"ancillaries"`
	Balances balance `json:"balances"`
}

type reservation struct{
	Flights []flight `json:"flights"`
	Passengers []passenger `json:"passengers"`
}

type reservationWithPNR struct{
	Reservation reservation `json:"reservation"`
	Pnr string `json:"pnr"`
}

type response struct{
	Status int `json:"status"`
	Message string `json:"message"`
	Data map[string]interface{} `json:"data"`
}

func ConnectDB() *mongo.Client {
    client, err := mongo.NewClient(options.Client().ApplyURI(CONNECTION_STRING))
    if err != nil {
        log.Fatal(err)
    }

    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    //ping the database
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB")
    return client
}

func getCollection(client *mongo.Client, collectionName string) *mongo.Collection{
	collection := client.Database(DB_NAME).Collection(collectionName)
	return collection
}

// func updateFlight() gin.HandlerFunc{
// 	return func(c *gin.Context){
// 		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//         var singleFlight flight
//         defer cancel()

// 		collection := getCollection(CLIENT, "flights")

// 		origin := c.Query("origen")
// 		destination := c.Query("destino")
// 		date := c.Query("fecha")
// 		number := c.Query("numero")

// 		if err := c.BindJSON(&singleFlight); err != nil {
//             c.JSON(http.StatusBadRequest, response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
//             return
//         }
// 		update :=bson.M{

// 			"number": singleFlight.Number,
// 			"origin": singleFlight.Origin,
// 			"destination": singleFlight.Destination,
// 			"departureTime": singleFlight.DepartureTime,
// 			"arrivalTime": singleFlight.ArrivalTime,
// 			"date": singleFlight.Date,
// 			"plane": singleFlight.Plane,
// 			"ancillaries": singleFlight.Ancillaries,
// 		}
// 		err := collection.UpdateOne(ctx, bson.M{"number": number, "origin": origin, "destination": destination, "date": date}, bson.M{"$set":update} )

// 		if err != nil {
//             c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
//             return
//         }


// 	}
// }

func getFlight() gin.HandlerFunc{
	return func(c *gin.Context){

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        var singleFlight flight
        defer cancel()

		collection := getCollection(CLIENT, "flights")

		origin := c.Query("origen")
		destination := c.Query("destino")
		date := c.Query("fecha")

		err := collection.FindOne(ctx, bson.M{"origen": origin, "destino": destination, "fecha": date}).Decode(&singleFlight)

		if err != nil {
            c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        c.JSON(http.StatusOK, response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": singleFlight}})

	}
}
