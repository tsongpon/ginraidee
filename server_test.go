package main_test

import (
	"database/sql"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/golang-migrate/migrate"
	"github.com/labstack/echo"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/tsongpon/ginraidee/adapter"
	"github.com/tsongpon/ginraidee/controller"
	"github.com/tsongpon/ginraidee/service"
	v1Controller "github.com/tsongpon/ginraidee/v1/controller"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"time"
)

const integrationPostgresqlContainerName = "integrationTest_postgresql"

var (
	postgresqlContainerId string
	dockerClient          *docker.Client
	integrationTestDB     *sql.DB
	postgresPort          string
)

func TestMain(m *testing.M) {
	setupTest()
	retCode := m.Run()
	tearDown()
	// require for notifying build tool if there is failed test
	os.Exit(retCode)
}

func TestPing(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/ping")

	pingController := controller.NewPingController(integrationTestDB)
	// Assertions
	if assert.NoError(t, pingController.Ping(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "pong", rec.Body.String())
	}
}

func TestScg(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/scg")

	scgController := controller.NewSCGController()
	// Assertions
	if assert.NoError(t, scgController.Echo(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "The Siam Cement PCL", rec.Body.String())
	}
}

func TestGetRestaurants(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, req.URL.String(), "/maps/api/geocode/json")
		// Send response to be tested
		rw.Write([]byte(`OK`))
	}))
	// Close the server when test finishes
	defer server.Close()

	e := echo.New()
	q := make(url.Values)
	q.Set("address", "Bangsue")
	req := httptest.NewRequest(http.MethodGet, "/v1/restaurants?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	placeAdapter := adapter.NewGooglePlaceAdapter()
	geoCodeAdapter := adapter.NewGoogleGeoCodeAdapter()
	lineAdapter := adapter.NewLineMessageAdapter()
	searchHistoryAdapter := adapter.NewSearchHistoryDBAdapter(integrationTestDB)
	ginRaiDeeService := service.NewGinRaiDeeService(placeAdapter, geoCodeAdapter, lineAdapter, searchHistoryAdapter)
	restaurantsController := v1Controller.NewRestaurantController(ginRaiDeeService)

	// Assertions
	expectedResponse := "{\"results\":[{\"placeId\":\"ChIJ1X_ASIGc4jARfsa8a3SI9wA\",\"name\":\"ร้านอาหาร ศิริชัยไก่ย่าง สาขาบางซื่อ (วงศ์สว่าง)\",\"link\":\"https://www.google.com/maps/place/?q=place_id:ChIJ1X_ASIGc4jARfsa8a3SI9wA\"}],\"next\":\"?address=Bangsue\\u0026pagetoken=gotonext\"}\n"
	if assert.NoError(t, restaurantsController.ListRestaurants(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
		assert.Equal(t, expectedResponse, rec.Body.String())
	}
}

//Util function
func setupTest() {
	go startMockServer()
	prepareDocker()
	prepareEnv()
	startContainer()
	//wait for container start properly
	time.Sleep(2 * time.Second)
	initialDBConnection()
}

func initialDBConnection() {
	dbHost := "localhost"
	dbPort := postgresPort
	dbUser := "postgres"
	dbPassword := "pingu123"
	dbName := "postgres"

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	m, err := migrate.New(
		"file://migrations", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	m.Steps(1)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}

	integrationTestDB = db
}

func prepareDocker() {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	dockerClient = client
}

func prepareEnv() {
	log.Println("preparing docker container(s) for test environment")
	freePort, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}
	postgresPort = strconv.Itoa(freePort)
	postgresContainerIdFromHost := getPostgresContainerId()
	if postgresContainerIdFromHost == "" {
		log.Println("pulling postgresql docker image(if missing), be patient")
		err := dockerClient.PullImage(docker.PullImageOptions{Repository: "postgres:11.5"}, docker.AuthConfiguration{})
		if err != nil {
			log.Fatalf("error while pulling image: %v", err.Error())
		}
		portBindings := map[docker.Port][]docker.PortBinding{
			"5432/tcp": {{HostPort: postgresPort}}}

		createContHostConfig := docker.HostConfig{
			PortBindings:    portBindings,
			PublishAllPorts: true,
			Privileged:      false}

		env := []string{"POSTGRES_PASSWORD=pingu123"}
		container, err := dockerClient.CreateContainer(docker.CreateContainerOptions{
			Name: integrationPostgresqlContainerName,
			Config: &docker.Config{
				Image: "postgres:11.5",
				Env:   env,
			},
			HostConfig: &createContHostConfig,
		})
		if err == nil {
			postgresqlContainerId = container.ID
		}
	} else {
		postgresqlContainerId = postgresContainerIdFromHost
	}
}

func getPostgresContainerId() string {
	cons, err := dockerClient.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		log.Fatalf("listing container error %v", err.Error())
	}
	for _, con := range cons {
		if con.Names[0] == "/"+integrationPostgresqlContainerName {
			return con.ID
		}
	}
	return ""
}

func startContainer() {
	portBindings := map[docker.Port][]docker.PortBinding{
		"5432/tcp": {{HostPort: postgresPort}}}

	err := dockerClient.StartContainer(postgresqlContainerId, &docker.HostConfig{PortBindings: portBindings})
	if err != nil {
		log.Printf("already running, keep going")
	}
}

func tearDown() {
	log.Printf("tearing down test enviroment")
	var err error
	err = dockerClient.StopContainer(postgresqlContainerId, 2)
	if err != nil {
		log.Printf("stop container error %v", err.Error())
	}
	err = dockerClient.RemoveContainer(docker.RemoveContainerOptions{ID: postgresqlContainerId, Force: true})
	if err != nil {
		log.Printf("remove container error %v", err.Error())
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func handleGeoCodeRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, googleGeoCodeResponse)
}

func handlePlaceRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, googlePlaceResponse)
}

func startMockServer() {
	http.HandleFunc("/maps/api/geocode/json", handleGeoCodeRequest)
	http.HandleFunc("/maps/api/place/nearbysearch/json", handlePlaceRequest)
	fmt.Println("starting mock server")
	http.ListenAndServe(":8080", nil)
}

const googleGeoCodeResponse = `{
   "results" : [
      {
         "address_components" : [
            {
               "long_name" : "Bang Sue",
               "short_name" : "Bang Sue",
               "types" : [ "political", "sublocality", "sublocality_level_1" ]
            },
            {
               "long_name" : "Bangkok",
               "short_name" : "Bangkok",
               "types" : [ "administrative_area_level_1", "political" ]
            },
            {
               "long_name" : "Thailand",
               "short_name" : "TH",
               "types" : [ "country", "political" ]
            }
         ],
         "formatted_address" : "Bang Sue, Bangkok, Thailand",
         "geometry" : {
            "bounds" : {
               "northeast" : {
                  "lat" : 13.8496853,
                  "lng" : 100.5449568
               },
               "southwest" : {
                  "lat" : 13.7972891,
                  "lng" : 100.5063032
               }
            },
            "location" : {
               "lat" : 13.828253,
               "lng" : 100.5284507
            },
            "location_type" : "APPROXIMATE",
            "viewport" : {
               "northeast" : {
                  "lat" : 13.8496853,
                  "lng" : 100.5449568
               },
               "southwest" : {
                  "lat" : 13.7972891,
                  "lng" : 100.5063032
               }
            }
         },
         "place_id" : "ChIJX5dpCoGb4jARUE_iXbIAAQM",
         "types" : [ "political", "sublocality", "sublocality_level_1" ]
      }
   ],
   "status" : "OK"
}`

const googlePlaceResponse = `{
  "html_attributions": [],
  "next_page_token": "gotonext",
  "results": [
    {
      "geometry": {
        "location": {
          "lat": 13.8306349,
          "lng": 100.5346448
        },
        "viewport": {
          "northeast": {
            "lat": 13.8319423802915,
            "lng": 100.5359315802915
          },
          "southwest": {
            "lat": 13.8292444197085,
            "lng": 100.5332336197085
          }
        }
      },
      "icon": "https://maps.gstatic.com/mapfiles/place_api/icons/restaurant-71.png",
      "id": "ab3a3efb5bcf0406036ae30d932f930bb58cdfbc",
      "name": "ร้านอาหาร ศิริชัยไก่ย่าง สาขาบางซื่อ (วงศ์สว่าง)",
      "opening_hours": {
        "open_now": false
      },
      "photos": [
        {
          "height": 612,
          "html_attributions": [
            "<a href=\"https://maps.google.com/maps/contrib/103155931138804094925/photos\">Tavicha Thongjirachote</a>"
          ],
          "photo_reference": "CmRaAAAAONzftnVTxcshInObwR4rUYDcShqrlJ5kT7kVnO-d6NqTbnSaGpLxY9UYHNHJ8P1E9diKSaRgJrPT7CoKxcI-K1PxYvgO9xnNLbfZuL3wp59SMxUjMxkqJQI4PBfseNCGEhDrygS51TrRmvH1TteaEFjrGhSbnFE7GGRxYokVdFBx-PhDDDxz0A",
          "width": 816
        }
      ],
      "place_id": "ChIJ1X_ASIGc4jARfsa8a3SI9wA",
      "plus_code": {
        "compound_code": "RGJM+7V Bangkok, Thailand",
        "global_code": "7P52RGJM+7V"
      },
      "price_level": 2,
      "rating": 4,
      "reference": "ChIJ1X_ASIGc4jARfsa8a3SI9wA",
      "scope": "GOOGLE",
      "types": [
        "restaurant",
        "food",
        "point_of_interest",
        "establishment"
      ],
      "user_ratings_total": 50,
      "vicinity": "59 Ratchadaphisek 37 Alley"
    }
  ],
  "status": "OK"
}`