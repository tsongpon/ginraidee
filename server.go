package main

import (
	"fmt"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tsongpon/ginraidee/adapter"
	"github.com/tsongpon/ginraidee/controller"
	"github.com/tsongpon/ginraidee/service"
	v1Controller "github.com/tsongpon/ginraidee/v1/controller"
	"log"
	"net/http"
	"os"

	"database/sql"
	_ "github.com/lib/pq"

	_ "github.com/golang-migrate/migrate/source/file"
)

func main() {
	log.Println("Staring server")

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "pingu123")
	dbName := getEnv("DB_NAME", "postgres")

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	m, err := migrate.New(
		"file://migrations", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	m.Steps(1)

	//psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	//	dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	ping := controller.NewPingController(db)
	scg := controller.NewSCGController()

	placeAdapter := adapter.NewGooglePlaceAdapter()
	geoCodeAdapter := adapter.NewGoogleGeoCodeAdapter()
	lineAdapter := adapter.NewLineMessageAdapter()
	searchHistoryAdapter := adapter.NewSearchHistoryDBAdapter(db)

	ginRaiDeeService := service.NewGinRaiDeeService(placeAdapter, geoCodeAdapter, lineAdapter, searchHistoryAdapter)

	lineHookController := v1Controller.NewLineHookController(ginRaiDeeService)

	restaurantController := v1Controller.NewRestaurantController(ginRaiDeeService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.GET("/scg", scg.Echo)
	e.GET("/ping", ping.Ping)

	e.GET("/v1/restaurants", restaurantController.ListRestaurants)
	e.POST("/v1/linehook", lineHookController.HandleMessage)
	e.Logger.Fatal(e.Start(":" + port))
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}