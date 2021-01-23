package web

import (
	"context"
	"k8sportal/model"
	"k8sportal/mongodb"
	"net/http"

	"github.com/foolin/goview/supports/ginview"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

//StartWebserver Gets entries from db and presents them via the HTTP endpoint
func StartWebserver(ctx context.Context, mongoClient *mongo.Client, mongodbDatabase string, mongodbCollection string) {

	router := gin.New()
	router.Use(gin.Recovery())

	router.HTMLRender = ginview.Default()

	router.GET("/services", func(ginCtx *gin.Context) {

		services := handleGetExistingServices(ginCtx, ctx, mongoClient, mongodbDatabase, mongodbCollection)

		ginCtx.HTML(http.StatusOK, "index", gin.H{
			"title":    "K8S Portal",
			"services": services,
		})
	})

	router.GET("/servicesapi", func(ginCtx *gin.Context) {
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

func handleGetExistingServices(ginCtx *gin.Context, ctx context.Context, mongoClient *mongo.Client, mongodbDatabase string, mongodbCollection string) []*model.Service {
	var loadedServices, err = mongodb.GetAllServices(ctx, mongoClient, mongodbDatabase, mongodbCollection)
	if err != nil {
		ginCtx.JSON(http.StatusNotFound, gin.H{"msg": err})
	}

	var existingServices []*model.Service

	for _, service := range loadedServices {
		if service.Exists() {
			existingServices = append(existingServices, service)
		}
	}

	return existingServices

}
