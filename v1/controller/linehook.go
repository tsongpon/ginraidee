package controller

import (
	"github.com/labstack/echo"
	"github.com/tsongpon/ginraidee/adapter"
	"github.com/tsongpon/ginraidee/model"
	"github.com/tsongpon/ginraidee/service"
	"github.com/tsongpon/ginraidee/v1/transport"
	"log"
	"net/http"
)

type LineHookController struct {
	service     *service.GinRaiDeeService
	lineAdapter *adapter.MessageAdapter
}

func NewLineHookController(service *service.GinRaiDeeService) *LineHookController {
	controller := new(LineHookController)
	controller.service = service
	return controller
}

func (c *LineHookController) HandleMessage(ctx echo.Context) error {
	log.Println("Start handle line message")
	eventTransport := transport.LineEventTransport{}
	if err := ctx.Bind(&eventTransport); err != nil {
		return err
	}
	lineEvent := model.LineEvent{}
	lineEvent.Type = eventTransport.Events[0].Type
	lineEvent.Message.Type = eventTransport.Events[0].Message.Type
	lineEvent.Message.Text = eventTransport.Events[0].Message.Text
	lineEvent.Timestamp = eventTransport.Events[0].Timestamp
	lineEvent.Source.Type = eventTransport.Events[0].Source.Type
	lineEvent.Source.UserID = eventTransport.Events[0].Source.UserID
	lineEvent.ReplyToken = eventTransport.Events[0].ReplyToken

	if err := c.service.HandleLineMessage(lineEvent); err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.String(http.StatusOK, "ok")
}
