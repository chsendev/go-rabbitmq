package rmq

import (
	"github.com/chsendev/go-rabbitmq/config"
	"time"
)

type Opt func(*config.RabbitmqConfig)

func WithUrl(url string) Opt {
	return func(conf *config.RabbitmqConfig) {
		conf.Url = url
	}
}

func WithPrefetch(size int) Opt {
	return func(conf *config.RabbitmqConfig) {
		conf.Listener.Prefetch = size
	}
}

func WithAckMode(mode config.AcknowledgeMode) Opt {
	return func(conf *config.RabbitmqConfig) {
		conf.Listener.AcknowledgeMode = mode
	}
}

func WithLogLevel(level string) Opt {
	return func(conf *config.RabbitmqConfig) {
		conf.Log.Level = level
	}
}

func WithRetry(initialInterval time.Duration, multiplier int, maxAttempts int) Opt {
	return func(conf *config.RabbitmqConfig) {
		conf.Retry.InitialInterval = initialInterval
		conf.Retry.Multiplier = multiplier
		conf.Retry.MaxAttempts = maxAttempts
	}
}
