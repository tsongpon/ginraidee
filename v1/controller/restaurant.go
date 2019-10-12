package controller

import (
	"github.com/labstack/echo"
	"net/http"
)

type RestaurantController struct {

}

func NewRestaurantController() *RestaurantController {
	return new(RestaurantController)
}

func (c *RestaurantController) Echo(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "The Siam Cement PCL")
}

