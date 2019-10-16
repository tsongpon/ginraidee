# GinRaiDee    
 Restaurant finder service which provide:    
    
 - Line chatbot integration    
 - RESTFul provide restaurant information    
  
## Technology stack  
 - Golang for programming language  
 - Postgresql for data store  
 - Docker for setup integration test environment    
    
## Line chatbot integration    
 Become friend with `@389weezk` and then you can just type in area name linebot will give some interesting restaurants on that area.    
    
## RESTFul provide restaurant information    
 Provide endpoint to list restaurant information by given address/area.    
    
 GET https://ginraidee.herokuapp.com/v1/restaurants?address=bangsue    
## Architecture    
 ![enter image description here](https://raw.githubusercontent.com/tsongpon/ginraidee/master/diagram/architecture.png)    
    
- **controller layer**: handle http request from outside    
      - linehook controller: handle webhook call made by line messaging api when user send message to chatbot    
      - restaurant controller: handle call  made by client to RESTFul api to get restaurant information     
- **service layer**: handle business logic    
- **adapter/repository layer**: make call to external service to retrieve information    
      - line message adapter: handle communication with line API    
      - google geocode adapter: ask for coordinate (lat, lon) by given area name    
      - google place api: ask for place with specific type(restaurant in this case) by   given coordinate    
      - database adapter: handle communication with data store to store user search history    
    
## Sequence diagrams    
 Sequence diagram show how chatbot service interact with line messaging API, Google place API and aatabase    
    
![enter image description here](https://raw.githubusercontent.com/tsongpon/ginraidee/master/diagram/sequence.png)

## Run service  
 Required environment variable:
 - DB_HOST: database host name
 - DB_NAME: database name
 - DB_PASSWORD: database password
 - DB_PORT: database port number
 - DB_USER: database user name
 - GEOCODE_ENDPOINT: google geocode API endpoint
 - GOOGLE_API_KEY: google API key
 - LINE_TOKEN: line messaging API token
 - PLACE_API_ENDPOINT: google place API endpoint

run command to start server:

    go run server.go

## Integration test

    export GEOCODE_ENDPOINT=http://localhost:8080/maps/api/geocode/json;export PLACE_API_ENDPOINT=http://localhost:8080/maps/api/place/nearbysearch/json;  go test -v