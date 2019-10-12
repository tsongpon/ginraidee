package controller

import (
	"github.com/labstack/echo"
	"net/http"
)

type SCGController struct {
}

func NewSCGController() *SCGController {
	return new(SCGController)
}

func (c *SCGController) Echo(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "The Siam Cement PCL")
}
