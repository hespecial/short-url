package response

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"short-url/internal/common/code"
)

type Response struct {
	Code    code.Code   `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Handle(ctx *gin.Context, c code.Code, data interface{}, message ...string) {
	msg := c.GetMessage()
	if len(message) > 0 {
		msg = message[0]
	}

	switch c {
	case code.Success:
		break
	case code.BadRequest:
		slog.Error(msg)
		msg = c.GetMessage()
	default:
		slog.Info(msg)
	}

	ctx.JSON(http.StatusOK, Response{
		Code:    c,
		Message: msg,
		Data:    data,
	})
}

func InvalidParams(ctx *gin.Context, err error) {
	Handle(ctx, code.InvalidParams, nil, err.Error())
}

func InvalidRequest(ctx *gin.Context, message ...string) {
	Handle(ctx, code.InvalidRequest, nil, message...)
}

func BadRequest(ctx *gin.Context, err error) {
	Handle(ctx, code.BadRequest, nil, err.Error())
}

func Unauthorized(ctx *gin.Context) {
	Handle(ctx, code.Unauthorized, nil)
}

func Success(ctx *gin.Context, data ...interface{}) {
	if len(data) == 0 {
		Handle(ctx, code.Success, struct{}{})
	} else {
		Handle(ctx, code.Success, data[0])
	}
}
