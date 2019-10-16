package controller

import (
	"database/sql"
	"github.com/labstack/echo"
	"net/http"
)

type PingController struct {
	db *sql.DB
}

func NewPingController(db *sql.DB) *PingController {
	pingController :=  new(PingController)
	pingController.db = db
	return pingController
}

func (c *PingController) Ping(ctx echo.Context) error {
	err := c.db.Ping()
	if err != nil {
		return ctx.String(http.StatusServiceUnavailable, err.Error())
	}
	return ctx.String(http.StatusOK, "pong")
}
