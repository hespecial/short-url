package controller

import (
	"net/http"
	"short-url/common/enum"
	"short-url/common/response"
	"short-url/internal/service"

	"github.com/gin-gonic/gin"
)

type ShortUrlController struct {
	svc service.ShortUrlService
}

func NewShortUrlController(svc service.ShortUrlService) *ShortUrlController {
	return &ShortUrlController{
		svc: svc,
	}
}

func (ctl *ShortUrlController) RevertToShortUrl(c *gin.Context) {
	var form struct {
		Url      string        `binding:"required" json:"url"`
		Priority enum.Priority `binding:"" json:"priority"`
		Comment  string        `binding:"" json:"comment"`
	}
	if err := c.ShouldBind(&form); err != nil {
		response.InvalidParams(c, err)
		return
	}

	shortUrl, err := ctl.svc.RevertToShortUrl(c, form.Url, form.Priority, form.Comment)
	if err != nil {
		response.InvalidRequest(c, err.Error())
		return
	}

	data := map[string]interface{}{
		"short_url": shortUrl,
	}
	response.Success(c, data)
}

func (ctl *ShortUrlController) RedirectToOriginalUrl(c *gin.Context) {
	shorUrlCode := c.Param("code")
	urlMapping, err := ctl.svc.GetUrlMappingByShortUrlCode(c, shorUrlCode)
	if err != nil {
		response.InvalidRequest(c, err.Error())
		return
	}

	if err = ctl.svc.LogAccess(c, urlMapping.Id, c.ClientIP(), c.GetHeader("User-Agent")); err != nil {
		response.BadRequest(c, err)
		return
	}

	c.Redirect(http.StatusFound, urlMapping.OriginalUrl)
}
