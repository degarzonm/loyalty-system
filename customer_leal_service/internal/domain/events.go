package domain

type EventProducer interface {
	SendMessage(topic string, message interface{}) error
}

type EventListener interface {
	Listen() error
}
