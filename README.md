# GinRaiDee

Restaurant finder service provide:

 - Line chatbot integration
 - RESTFul provide restaurant information

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
		- google place api: ask for place with specific type(restaurant in this case) by 	given coordinate
		- database adapter: handle communication with data store to store user search history

## Sequence diagrams

Sequence diagram show how chatbot service interact with line messaging API, Google place API and aatabase

```mermaid
sequenceDiagram
Line Messing API ->> LineHook Controller: message event
LineHook Controller ->> GinRaiDee Service: HandleLineMessage()
GinRaiDee Service ->> Geocode Adapter: GetLocation("bangsue")
Geocode Adapter ->> GinRaiDee Service: location(lat, lon)
GinRaiDee Service ->> Plcae Adapter: GetPlaces("restaurant", lat, lon)
Plcae Adapter ->> GinRaiDee Service: list of restaurant
GinRaiDee Service ->> Database Adapter: Save(userID, search keyword)
Database Adapter ->> Postgresql: INSERT INTO searchhistory
Database Adapter ->> GinRaiDee Service: Saved search hisrtory
GinRaiDee Service ->> LineHook Controller: 
LineHook Controller ->> Line Messing API: http status OK
```