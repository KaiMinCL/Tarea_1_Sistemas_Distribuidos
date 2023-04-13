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
	//"encoding/json"
)

const(
	DB_NAME = "aerolinea"
)

type plane struct {
	modelo string `bson:"modelo"`
	numero_de_serie string `bson:"numero_de_serie"`
	stock_de_pasajeros int `bson:"stock_de_pasajeros"`
}

type ancillary struct {
	nombre string `bson:"nombre"`
	stock int `bson:"stock"`
	ssr string `bson:"ssr"`
}

type passenger_ancillary struct{
	DepartureFlight []ancillary `bson:"departureFlight"`
	ReturnFlight []ancillary `bson:"returnFlight"`
}

type flight struct{
	id primitive.ObjectID `bson:"_id"`
	numero_vuelo string `bson:"numero_vuelo"`
	origen string `bson:"origen"`
	destino string `bson:"destino"`
	hora_salida string `bson:"hora_salida"`
	hora_llegada string `bson:"hora_llegada"`
	fecha string `bson:"fecha"`
	avion plane `bson:"avion"`
	ancillaries []ancillary `bson:"ancillaries"`
}

type balance struct{
	AncillariesDepartureFlight int `bson:"ancillariesDepartureFlight"`
	DepartureFlight int `bson:"departureFlight"`
	AncillariesReturnFlight int `bson:"ancillariesReturnFlight"`
	ReturnFlight int `bson:"returnFlight"`
}

type passenger struct{
	FirstName string `bson:"firstName"`
	Surname string `bson:"surname"`
	Age int `bson:"age"`
	Ancillaries passenger_ancillary `bson:"ancillaries"`
	Balances balance `bson:"balances"`
}

type reservation struct{
	Flights []flight `bson:"flights"`
	Passengers []passenger `bson:"passengers"`
}

type reservationWithPNR struct{
	Reservation reservation `bson:"reservation"`
	Pnr string `bson:"pnr"`
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
