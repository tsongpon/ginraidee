package controller

import (
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo"
	"github.com/tsongpon/ginraidee/adapter"
	"github.com/tsongpon/ginraidee/v1/transport"
	"log"
	"net/http"
	"os"
	"strconv"
)

var accessToken = os.Getenv("LINE_TOKEN")
var lineReplyEndpoint = "https://api.line.me/v2/bot/message/reply"

type LineHookController struct {
}

func NewLineHookController() *LineHookController {
	return new(LineHookController)
}

func (c *LineHookController) HandleMessage(ctx echo.Context) error {
	log.Println("Start handle line message")
	event := transport.LineEvent{}
	if err := ctx.Bind(&event); err != nil {
		return err
	}

	placeAdapter := adapter.NewGooglePlaceAdapter();
	places := placeAdapter.GetPlaces("restaurant", 13.828253, 100.5284507)
	replyMessage := ""
	for _, each := range places {
		replyMessage = replyMessage + each.Name + " (" + strconv.Itoa(int(each.Ratting)) + ")" + "\n"
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
	log.Printf("response status %d", rep.StatusCode())

	return ctx.String(http.StatusOK, "ok")
}
