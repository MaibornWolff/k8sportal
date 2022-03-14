package web

import (
	"k8sportal/k8sclient"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/foolin/goview/supports/ginview"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

)

//StartWebserver Gets entries from db and presents them via the HTTP endpoint
func StartWebserver() {

	router := gin.New()
	router.Use(cors.Default())
	router.Use(gin.Recovery())

	router.HTMLRender = ginview.Default()

	router.GET("/ui/services", func(ginCtx *gin.Context) {
		log.Info().Msg("calling ui")
		services := k8sclient.GetAllServices()
		log.Info().Msgf("services: %v", services)
		ginCtx.HTML(http.StatusOK, "index", gin.H{
			"title":    "K8S Portal",
			"services": services,
		})
	})

	router.GET("/api/services", func(ginCtx *gin.Context) {
		handleGetServices(ginCtx)
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
