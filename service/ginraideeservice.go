package service

import (
	"github.com/tsongpon/ginraidee/adapter"
	"github.com/tsongpon/ginraidee/model"
	"github.com/tsongpon/ginraidee/v1/transport"
	"strconv"
)

type GinRaiDeeService struct {
	placeAdapter   adapter.PlaceAdapter
	geoCodeAdapter adapter.GeoCodeAdapter
	lineAdapter adapter.MessageAdapter
}

func NewGinRaiDeeService(placeAdapter adapter.PlaceAdapter,
	geoCodeAdapter adapter.GeoCodeAdapter, lineAdapter adapter.MessageAdapter) *GinRaiDeeService {
	service := new(GinRaiDeeService)
	service.placeAdapter = placeAdapter
	service.geoCodeAdapter = geoCodeAdapter
	service.lineAdapter = lineAdapter
	return service
}

func (s *GinRaiDeeService) GetRestaurants(lineEvent model.LineEvent) error {
	var err error
	location, err := s.geoCodeAdapter.GetLocation(lineEvent.Message.Text)
	if err != nil {
		return err
	}
	restaurants, err :=  s.placeAdapter.GetPlaces("restaurant", location.Lat, location.Lng)

	replyMessage := ""
	for _, each := range restaurants {
		replyMessage = replyMessage + each.Name + " (" + strconv.Itoa(int(each.Ratting)) + ")" + "\n"
		//replyMessage = replyMessage + each.MapLink + "\n\n"
	}

	reply := transport.LineReply{}
	reply.ReplyToken = lineEvent.ReplyToken
	message := transport.ReplyMessage{}
	message.Text = truncateString(replyMessage, 2000)
	message.Type = "text"
	reply.Messages = []transport.ReplyMessage{message}

	if err := s.lineAdapter.Reply(reply); err != nil {
		return err
	}

	return nil
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}