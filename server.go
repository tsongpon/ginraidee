package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tsongpon/ginraidee/controller"
	v1Controller "github.com/tsongpon/ginraidee/v1/controller"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Starting server")
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	ping := controller.NewPingController()
	scg := controller.NewSCGController()
	lineHookController := v1Controller.NewLineHookController()

	e.GET("/scg", scg.Echo)
	e.GET("/ping", ping.Ping)

	e.POST("/v1/linehook", lineHookController.HandleMessage)
	e.Logger.Fatal(e.Start(":" + port))
}
