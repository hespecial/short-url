package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"short-url/internal/common/enum"
	"short-url/internal/common/response"
	"short-url/internal/service"
)

type ShortUrlController struct {
	svc *service.ShortUrlService
}

func NewShortUrlController(svc *service.ShortUrlService) *ShortUrlController {
	return &ShortUrlController{
		svc: svc,
	}
}

func (hdl *ShortUrlController) RevertToShortUrl(c *gin.Context) {
	var form struct {
		Url      string        `binding:"required" json:"url"`
		Priority enum.Priority `binding:"" json:"priority"`
		Comment  string        `binding:"" json:"comment"`
	}
	if err := c.ShouldBind(&form); err != nil {
		response.InvalidParams(c, err)
		return
	}

	shortUrl, err := hdl.svc.RevertToShortUrl(c, form.Url, form.Priority, form.Comment)
	if err != nil {
		response.InvalidRequest(c, err.Error())
		return
	}

	data := map[string]interface{}{
		"short_url": shortUrl,
	}
	response.Success(c, data)
}

func (hdl *ShortUrlController) RedirectToOriginalUrl(c *gin.Context) {
	shorUrlCode := c.Param("code")
	urlMapping, err := hdl.svc.GetUrlMappingByShortUrlCode(c, shorUrlCode)
	if err != nil {
		response.InvalidRequest(c, err.Error())
		return
	}

	if err = hdl.svc.ProcessAccess(c, urlMapping.Id, c.ClientIP()); err != nil {
		response.BadRequest(c, err)
		return
	}

	if err = hdl.svc.LogAccess(c, urlMapping.Id, c.ClientIP(), c.GetHeader("User-Agent")); err != nil {
		response.BadRequest(c, err)
		return
	}

	c.Redirect(http.StatusFound, urlMapping.OriginalUrl)
}
