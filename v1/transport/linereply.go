package transport

type LineReply struct {
	ReplyToken string         `json:"replyToken"`
	Messages   []ReplyMessage `json:"messages"`
}

type ReplyMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
