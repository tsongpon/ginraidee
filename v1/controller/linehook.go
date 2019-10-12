package controller

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo"
	"github.com/tsongpon/ginraidee/v1/transport"
	"net/http"
)

type LineHookController struct {
}

func NewLineHookController() *LineHookController {
	return new(LineHookController)
}

func (c *LineHookController) HandleMessage(ctx echo.Context) error {
	event := transport.LineEvent{}
	if err := ctx.Bind(&event); err != nil {
		return err
	}
	reply := transport.LineReply{}
	reply.ReplyToken = event.Events[0].ReplyToken
	message := transport.ReplyMessage{}
	message.Text = "Kin Rai Dee"
	message.Type = "text"
	reply.Messages = []transport.ReplyMessage{message}

	client := resty.New()
	_, err := client.R().
		SetBody(reply).
		SetAuthToken(event.Events[0].ReplyToken).
		SetHeader("Content-Type", "application/json").
		Post("https://api.line.me/v2/bot/message")

	if err != nil {
		fmt.Print(err.Error())
	}

	return ctx.String(http.StatusOK, "ok")
}
