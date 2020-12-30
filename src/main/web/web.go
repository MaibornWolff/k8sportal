package web

import (
	"context"
	"k8sportal/mongodb"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

//StartWebserver Gets entries from db and presents them via the HTTP endpoint
func StartWebserver(ctx context.Context, mongoClient *mongo.Client, mongodbDatabase string, mongodbCollection string) {

	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/services", func(ginCtx *gin.Context) {
		handleGetServices(ginCtx, ctx, mongoClient, mongodbDatabase, mongodbCollection)
	})
	router.Run(":80")
}

func handleGetServices(ginCtx *gin.Context, ctx context.Context, mongoClient *mongo.Client, mongodbDatabase string, mongodbCollection string) {
	var loadedServices, err = mongodb.GetAllServices(ctx, mongoClient, mongodbDatabase, mongodbCollection)
	if err != nil {

		ginCtx.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}

	ginCtx.JSON(http.StatusOK, loadedServices)

}
