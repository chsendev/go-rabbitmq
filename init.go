package rmq

import (
	"github.com/cscoder0/go-rabbitmq/config"
	"github.com/cscoder0/go-rabbitmq/log"
)

func Init(conf *config.RabbitmqConfig) {
	config.Init(conf)
	log.Init()
}
