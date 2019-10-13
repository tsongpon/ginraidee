package adapter

type MessageAdapter interface {
	Reply(interface{}) error
}
