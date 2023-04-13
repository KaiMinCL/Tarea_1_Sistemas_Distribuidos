package main


import (
	"log"
    "os"
    "github.com/joho/godotenv"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
	"context"
	"time"
)

var(
	SERVER string
	PORT string
	CONNECTION_STRING string
	CLIENT *mongo.Client
)

func EnvMongoURI(){
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	SERVER = os.Getenv("SERVER")
	PORT = os.Getenv("PORT")
	CONNECTION_STRING = os.Getenv("CONNECTION_STRING")
}


func Routes(router *gin.Engine){
	router.GET("/api/vuelo", getFlight())
	//router.PUT("/api/vuelo", updateFlight())
}

func setupData(){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defaultFlight := flight{
		id:primitive.NewObjectID(),
		numero_vuelo:"323",
		origen:"IQQ",
		destino:"SCL",
		hora_salida:"20h38",
		hora_llegada:"22H50",
		fecha:"14/04/2023",
		avion:plane{modelo:"A320neo", numero_de_serie: "12345", stock_de_pasajeros:90},
		ancillaries:[]ancillary{{nombre:"Equipaje_de_mano", stock:68, ssr:"BGH"},
			{nombre:"Equipaje_de_bodega", stock:92, ssr:"BGR"},
			{nombre:"Asiento", stock:90, ssr:"STDF"},
			{nombre:"Embarque y Check In prioritario", stock:79, ssr:"PAXS"},
			{nombre:"Mascota en cabina", stock:4, ssr:"PTCR"},
			{nombre:"Mascota en bodega", stock:12, ssr:"AVIH"},
			{nombre:"Equipaje especial", stock:71, ssr:"SPML"},
			{nombre:"Acceso a Salón VIP", stock:36, ssr:"LNGE"},
			{nombre:"Wi-Fi a bordo", stock:57, ssr:"WIFI"},
		},
	}

	fmt.Println(defaultFlight)


	collection := getCollection(CLIENT, "flights")

	result, err := collection.InsertOne(ctx, defaultFlight)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(result)
}

func main(){
	router := gin.Default()

	EnvMongoURI()

	fmt.Println("---->", SERVER, PORT, CONNECTION_STRING)

	CLIENT = ConnectDB()
	Routes(router)
	router.Run(SERVER+":"+PORT);

}
