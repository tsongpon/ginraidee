package controller

import (
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo"
	"github.com/tsongpon/ginraidee/service"
	"github.com/tsongpon/ginraidee/v1/transport"
	"log"
	"net/http"
	"os"
	"strconv"
)

var accessToken = os.Getenv("LINE_TOKEN")
var lineReplyEndpoint = "https://api.line.me/v2/bot/message/reply"

type LineHookController struct {
	service *service.GinRaiDeeService
}

func NewLineHookController(service *service.GinRaiDeeService) *LineHookController {
	controller := new(LineHookController)
	controller.service = service
	return controller
}

func (c *LineHookController) HandleMessage(ctx echo.Context) error {
	log.Println("Start handle line message")
	var err error
	event := transport.LineEvent{}
	if err := ctx.Bind(&event); err != nil {
		return err
	}

	address := event.Events[0].Message.Text
	places, err := c.service.GetRestaurants(address)
	replyMessage := ""
	for _, each := range places {
		replyMessage = replyMessage + each.Name + " (" + strconv.Itoa(int(each.Ratting)) + ")" + "\n"
		replyMessage = replyMessage + each.MapLink + "\n\n"
	}

	reply := transport.LineReply{}
	reply.ReplyToken = event.Events[0].ReplyToken
	message := transport.ReplyMessage{}
	message.Text = replyMessage
	message.Type = "text"
	reply.Messages = []transport.ReplyMessage{message}

	client := resty.New()
	rep, err := client.R().
		SetBody(reply).
		SetAuthToken(accessToken).
		SetHeader("Content-Type", "application/json").
		Post(lineReplyEndpoint)

	if err != nil {
		log.Printf("error from line %s", err.Error())
	}
	log.Printf("response status %d, msg %s", rep.StatusCode(), rep.String())

	return ctx.String(http.StatusOK, "ok")
}
