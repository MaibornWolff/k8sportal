package web

import (
	"log"

	"github.com/gin-gonic/gin"
)

func StartWebserver() {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/lastChange", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	go r.Run(":80")

	log.Print("Webserver Works")

}
