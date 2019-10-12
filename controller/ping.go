package controller

import (
	"github.com/labstack/echo"
	"github.com/tsongpon/ginraidee/adapter"
	"log"
	"net/http"
)

type PingController struct {
}

func NewPingController() *PingController {
	return new(PingController)
}

func (c *PingController) Ping(ctx echo.Context) error {
	placeAdapter := adapter.NewGooglePlaceAdapter()
	places := placeAdapter.GetPlaces("restaurant", 13.828253, 100.5284507)
	log.Printf("places %v", places)
	return ctx.String(http.StatusOK, "pong")
}
