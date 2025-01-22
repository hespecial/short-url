package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"short-url/global"
	"short-url/internal/controller"
	"short-url/internal/repo"
	"short-url/internal/repo/cache"
	"short-url/internal/repo/dao"
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

	db := global.Db
	rdb := global.Redis
	svc := service.NewShortUrlService(repo.NewShortUrlRepo(
		dao.NewShortUrlDao(db),
		cache.NewShortUrlCache(rdb),
	))

	ctl := controller.NewShortUrlController(svc)
	{
		router.GET("/:code", ctl.RedirectToOriginalUrl)

		api.POST("/shorten", ctl.RevertToShortUrl)

		api.GET("/short-url/list")
		api.GET("/short-url/:id")
		api.PUT("/short-url/:id")
		api.DELETE("/short-url/:id")
	}

	log.Fatal(router.Run(fmt.Sprintf("%s:%d", global.Conf.App.Host, global.Conf.App.Port)))
}
