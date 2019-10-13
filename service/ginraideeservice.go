package service

import (
	"github.com/tsongpon/ginraidee/adapter"
	"github.com/tsongpon/ginraidee/model"
	"github.com/tsongpon/ginraidee/v1/transport"
	"strconv"
)

const maximumLineMessageLength = 2000

type GinRaiDeeService struct {
	placeAdapter   adapter.PlaceAdapter
	geoCodeAdapter adapter.GeoCodeAdapter
	lineAdapter    adapter.MessageAdapter
}

func NewGinRaiDeeService(placeAdapter adapter.PlaceAdapter,
	geoCodeAdapter adapter.GeoCodeAdapter, lineAdapter adapter.MessageAdapter) *GinRaiDeeService {
	service := new(GinRaiDeeService)
	service.placeAdapter = placeAdapter
	service.geoCodeAdapter = geoCodeAdapter
	service.lineAdapter = lineAdapter
	return service
}

func (s *GinRaiDeeService) HandleLineMessage(lineEvent model.LineEvent) error {
	var err error
	location, err := s.geoCodeAdapter.GetLocation(lineEvent.Message.Text)
	if err != nil {
		return err
	}
	restaurants, _, err := s.placeAdapter.GetPlaces("restaurant", location.Lat, location.Lng, "")

	var replyMessage string
	for _, each := range restaurants {

		var rating string
		if each.Rating > 0.0 {
			rating = " (" + strconv.Itoa(int(each.Rating)) + ")"
		}
		messageToAdd := each.Name + rating + "\n"
		if len(messageToAdd)+len(replyMessage) <= maximumLineMessageLength {
			replyMessage = replyMessage + messageToAdd
		} else {
			break
		}
		//replyMessage = replyMessage + each.MapLink + "\n\n"
	}

	reply := transport.LineReply{}
	reply.ReplyToken = lineEvent.ReplyToken
	message := transport.ReplyMessage{}
	message.Text = replyMessage
	message.Type = "text"
	reply.Messages = []transport.ReplyMessage{message}

	if err := s.lineAdapter.Reply(reply); err != nil {
		return err
	}

	return nil
}

func (s *GinRaiDeeService) GetRestaurants(address string, pageToken string) ([]model.Place, string, error) {
	var err error
	location, err := s.geoCodeAdapter.GetLocation(address)
	if err != nil {
		return nil, "", err
	}
	restaurants, pageToken, err := s.placeAdapter.GetPlaces("restaurant", location.Lat, location.Lng, pageToken)

	return restaurants, pageToken, nil
}
