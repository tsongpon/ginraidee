package controller

import (
	"github.com/labstack/echo"
	"github.com/tsongpon/ginraidee/service"
	"github.com/tsongpon/ginraidee/v1/transport"
	v1Transport "github.com/tsongpon/ginraidee/v1/transport"
	"net/http"
)

type RestaurantController struct {
	service *service.GinRaiDeeService
}

func NewRestaurantController(service *service.GinRaiDeeService) *RestaurantController {
	controller := new(RestaurantController)
	controller.service = service
	return controller
}

func (c *RestaurantController) Echo(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "The Siam Cement PCL")
}

func (c *RestaurantController) ListRestaurants(ctx echo.Context) error {
	address := ctx.QueryParam("address")
	pageToken := ctx.QueryParam("pagetoken")
	places, pageToken, err := c.service.GetRestaurants(address, pageToken)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	var transports []transport.PlaceTransport
	for _, each := range places {
		t := v1Transport.PlaceTransport{
			PlaceID: each.PlaceID,
			Name: each.Name,
			Link: each.MapLink,
		}
		transports = append(transports, t)
	}

	var nextPage *string
	nextStr := ctx.Path() + "?address=" + address + "&pagetoken=" + pageToken
	if len(pageToken) > 0 {
		nextPage = &nextStr
	}
	responseTransport := v1Transport.PlaceResponseTransport{
		Result:   transports,
		NextPage: nextPage,
	}
	return ctx.JSON(http.StatusOK, responseTransport)
}
