package controller

import (
	"github.com/labstack/echo"
	"net/http"
)

type PingController struct {

}

func NewPingController() *PingController {
	return new(PingController)
}

func (c *PingController) Ping(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "pong")
}
