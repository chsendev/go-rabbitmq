package rmq

import (
	"github.com/ChsenDev/go-rabbitmq/config"
	"github.com/ChsenDev/go-rabbitmq/log"
	"sync"
)

var configOnce sync.Once

func Init(url string, opt ...Opt) {
	configOnce.Do(func() {
		config.Conf.Url = url
		for _, o := range opt {
			o(config.Conf)
		}
		log.Init()
	})
}
