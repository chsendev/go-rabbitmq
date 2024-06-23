package mq

type ExchangeType string

const (
	Default ExchangeType = ""
	Fanout  ExchangeType = "fanout"
	Direct  ExchangeType = "direct"
	Topic   ExchangeType = "topic"
)
