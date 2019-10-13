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
	geoAdapter := adapter.NewGoogleGeoCodeAdapter()
	location, _ := geoAdapter.GetLocation("บางใหญ่")
	log.Printf("lat: %f , lng: %f", location.Lat, location.Lng)
	return ctx.String(http.StatusOK, "pong")
}
