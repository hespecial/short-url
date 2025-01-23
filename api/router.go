package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"short-url/api/handler"
	"short-url/global"
	"short-url/internal/repo"
	"short-url/internal/service"
)

func StartServer() {
	gin.SetMode(global.Conf.App.Mode)
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api := router.Group("/api")

	svc := service.NewShortUrlService(repo.NewShortUrlRepo(
		global.Db,
		global.Redis,
	))

	hdl := handler.NewShortUrlController(svc)
	{
		router.GET("/:code", hdl.RedirectToOriginalUrl)

		api.POST("/shorten", hdl.RevertToShortUrl)

		api.GET("/short-url/list")
		api.GET("/short-url/:id")
		api.PUT("/short-url/:id")
		api.DELETE("/short-url/:id")
	}

	log.Fatal(router.Run(fmt.Sprintf("%s:%d", global.Conf.App.Host, global.Conf.App.Port)))
}
