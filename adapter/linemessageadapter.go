package adapter

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/tsongpon/ginraidee/v1/transport"
	"log"
	"os"
)

var accessToken = os.Getenv("LINE_TOKEN")
var lineReplyEndpoint = "https://api.line.me/v2/bot/message/reply"

type LineMessageAdapter struct {
}

func NewLineMessageAdapter() *LineMessageAdapter {
	return new(LineMessageAdapter)
}

func (m *LineMessageAdapter) Reply(message interface{}) error {
	lineReplyMessage, ok := message.(transport.LineReply)
	if !ok {
		return errors.New("message is not line reply message")
	}
	client := resty.New()
	rep, err := client.R().
		SetBody(lineReplyMessage).
		SetAuthToken(accessToken).
		SetHeader("Content-Type", "application/json").
		Post(lineReplyEndpoint)

	if err != nil {
		log.Printf("error from line %s", err.Error())
	}
	log.Printf("response status %d, msg %s", rep.StatusCode(), rep.String())

	return nil
}
