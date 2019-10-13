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

//const (
//	host     = "localhost"
//	dbPort     = 5432
//	user     = "postgres"
//	password = "pingu123"
//	dbname   = "postgres"
//)

func main() {
	log.Println("Staring server")
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	//psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	//	host, dbPort, user, password, dbname)
	//
	//log.Println(psqlInfo)
	//db, err := sql.Open("postgres", psqlInfo)
	//if err != nil {
	//	panic(err)
	//}
	//defer db.Close()
	//
	//err = db.Ping()
	//if err != nil {
	//	panic(err)
	//}

	ping := controller.NewPingController()
	scg := controller.NewSCGController()

	placeAdapter := adapter.NewGooglePlaceAdapter()
	geoCodeAdapter := adapter.NewGoogleGeoCodeAdapter()
	lineAdapter := adapter.NewLineMessageAdapter()
	ginRaiDeeService := service.NewGinRaiDeeService(placeAdapter, geoCodeAdapter, lineAdapter)
	lineHookController := v1Controller.NewLineHookController(ginRaiDeeService)

	restaurantContoller := v1Controller.NewRestaurantController(ginRaiDeeService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.GET("/scg", scg.Echo)
	e.GET("/ping", ping.Ping)

	e.GET("/v1/restaurants", restaurantContoller.ListRestaurants)
	e.POST("/v1/linehook", lineHookController.HandleMessage)
	e.Logger.Fatal(e.Start(":" + port))
}
