package web

import (
	"k8sportal/mongodb"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var mongodbdatabase = "k8sportal"
var mongodbcollection = "portal-services"

//StartWebserver Gets entries from db and presents them
func StartWebserver(mongoClient *mongo.Client) {

	router := gin.Default()
	router.GET("/services", func(ctx *gin.Context) {
		handleGetServices(ctx, mongoClient)
	})
	go router.Run(":80")

	log.Print("Webserver Works")

}

func handleGetServices(c *gin.Context, mongoClient *mongo.Client) {
	var loadedServices, err = mongodb.GetAllServices(mongoClient)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"services": loadedServices})
}
