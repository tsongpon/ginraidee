package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tsongpon/ginraidee/controller"
	"net/http"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	ping := controller.NewPingController()
	scg := controller.NewSCGController()

	e.GET("/scg", scg.Echo)
	e.GET("/ping", ping.Ping)

	e.Logger.Fatal(e.Start(":5000"))
}