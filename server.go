package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tsongpon/ginraidee/adapter"
	"github.com/tsongpon/ginraidee/controller"
	"github.com/tsongpon/ginraidee/service"
	v1Controller "github.com/tsongpon/ginraidee/v1/controller"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Staring server")
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	ping := controller.NewPingController()
	scg := controller.NewSCGController()

	placeAdapter := adapter.NewGooglePlaceAdapter()
	geoCodeAdapter := adapter.NewGoogleGeoCodeAdapter()
	ginRaiDeeService := service.NewGinRaiDeeService(placeAdapter, geoCodeAdapter)
	lineHookController := v1Controller.NewLineHookController(ginRaiDeeService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.GET("/scg", scg.Echo)
	e.GET("/ping", ping.Ping)

	e.POST("/v1/linehook", lineHookController.HandleMessage)
	e.Logger.Fatal(e.Start(":" + port))
}
