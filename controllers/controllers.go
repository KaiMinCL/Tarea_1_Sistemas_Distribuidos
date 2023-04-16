package controllers

import (
	"bd_aerolinea/models"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

/*func GetVuelos(c *gin.Context) {

	//vuelo := models.GetVuelos()

	if vuelo == nil || len(vuelo) == 0 {

		c.AbortWithStatus(http.StatusNotFound)

	} else {

		c.IndentedJSON(http.StatusOK, vuelo)

	}
}*/

func GetVuelo(c *gin.Context) {

	origenVuelo := c.Query("origen")
	destinoVuelo := c.Query("destino")
	fechaVuelo := c.Query("fecha")

	fmt.Println(origenVuelo, destinoVuelo, fechaVuelo)
	vuelo, err := models.GetVuelo(origenVuelo, destinoVuelo, fechaVuelo)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)

	} else {

		c.IndentedJSON(http.StatusOK, vuelo)

	}

}

func CreateVuelo(c *gin.Context) {

	var vuelo models.Flight

	if err := c.BindJSON(&vuelo); err != nil {

		c.AbortWithStatus(http.StatusBadRequest)
	} else {

		models.CreateVuelo(vuelo)
	}
}

func UpdateVuelo(c *gin.Context) {

	//var vuelo models.Flight

	numeroVuelo := c.Query("numero_vuelo")
	origenVuelo := c.Query("origen")
	destinoVuelo := c.Query("destino")
	fechaVuelo := c.Query("fecha")

	var updateBSON interface{}

	requestBody, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		// Handle error
	}

	err = bson.UnmarshalExtJSON(requestBody, true, &updateBSON)
	if err != nil {
		// Handle error
	}

	models.UpdateVuelo(numeroVuelo, origenVuelo, destinoVuelo, fechaVuelo, updateBSON)
	//c.IndentedJSON(http.StatusCreated, gin.H{"id_producto": vuelo.id_vuelo})

}

func DeleteVuelo(c *gin.Context) {

	numeroVuelo := c.Query("numero_vuelo")
	origenVuelo := c.Query("origen")
	destinoVuelo := c.Query("destino")
	fechaVuelo := c.Query("fecha")

	/*count :=*/
	models.DeleteVuelo(numeroVuelo, origenVuelo, destinoVuelo, fechaVuelo)

	/*if count > 0 {
		c.IndentedJSON(http.StatusAccepted, gin.H{"id_vuelo": id_vuelo})
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
	*/
}
