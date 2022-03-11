package web

import (
	"k8sportal/k8sclient"
	"k8sportal/model"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/foolin/goview/supports/ginview"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	v1 "k8s.io/api/core/v1"
)

//StartWebserver Gets entries from db and presents them via the HTTP endpoint
func StartWebserver() {

	router := gin.New()
	router.Use(cors.Default())
	router.Use(gin.Recovery())

	router.HTMLRender = ginview.Default()

	router.GET("/ui/services", func(ginCtx *gin.Context) {
		log.Info().Msg("calling ui")
		services := handleGetExistingServices(ginCtx)
		log.Info().Msgf("services: %v", services)
		ginCtx.HTML(http.StatusOK, "index", gin.H{
			"title":    "K8S Portal",
			"services": services,
		})
	})

	router.GET("/api/services", func(ginCtx *gin.Context) {
		handleGetServices(ginCtx)
	})

    router.GET("/api/ingresses", func(ginCtx *gin.Context) {
		handleGetIngresses(ginCtx)
	})

	router.Run(":80")
}

func handleGetServices(ginCtx *gin.Context) {
	//var loadedServices, err = mongodb.GetAllServices(ctx, mongoClient, mongodbDatabase, mongodbCollection)
	// if err != nil {

	// 	ginCtx.JSON(http.StatusNotFound, gin.H{"msg": err})
	// 	return
	// }
	ginCtx.JSON(http.StatusOK, k8sclient.GetAllServices())

}

func handleGetExistingServices(ginCtx *gin.Context) []*model.Service {
	// var loadedServices, err = mongodb.GetAllServices(ctx, mongoClient, mongodbDatabase, mongodbCollection)
	// if err != nil {
	// 	ginCtx.JSON(http.StatusNotFound, gin.H{"msg": err})
	// }

	var existingServices []*model.Service

	// TODO where to fill do the mapping?

	list := k8sclient.GetAllServices()
	for _, entry := range list {
		k8sservice := entry.(*v1.Service)
		service := &model.Service{
			ServiceName:   k8sservice.Name,
			Category:      "bla",
			ServiceExists: false,
			IngressRules:  []model.IngressRule{},
			IngressExists: false,
		}

		existingServices = append(existingServices, service)
	}
	// for _, service := range loadedServices {
	// 	if service.Exists() {
	// 		existingServices = append(existingServices, service)
	// 	}
	// }

	return existingServices

}

func handleGetIngresses(ginCtx *gin.Context) {
	ginCtx.JSON(http.StatusOK, k8sclient.GetAllIngresses())

}
