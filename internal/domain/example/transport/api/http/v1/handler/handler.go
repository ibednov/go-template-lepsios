package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpresponse "github.com/ibednov/go-lepsios/httpx/response"
	"github.com/ibednov/go-lepsios/log"
	"github.com/ibednov/go-template-lepsios/internal/domain/example/service"
	"github.com/ibednov/go-template-lepsios/internal/domain/example/transport/api/http/v1/response"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.GET("/hello", h.Hello)
}

func (h *Handler) Hello(c *gin.Context) {
	msg, err := h.svc.Hello(c.Request.Context())
	if err != nil {
		log.ErrorCtx(c.Request.Context(), "example.hello.failed", "error", err)
		httpresponse.Internal(c, "internal error")
		return
	}
	c.JSON(http.StatusOK, response.FromMessage(msg))
}
