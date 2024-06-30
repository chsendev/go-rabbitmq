package config

import "time"

type RabbitmqConfig struct {
	Url      string    `mapstructure:"url"`
	Log      *log      `mapstructure:"log"`
	Retry    *retry    `mapstructure:"retry"`
	Listener *listener `mapstructure:"listener"`
}

type AcknowledgeMode string

const (
	AcknowledgeModeNone   AcknowledgeMode = "none"
	AcknowledgeModeAuto   AcknowledgeMode = "auto"
	AcknowledgeModeManual AcknowledgeMode = "manual"
)

type log struct {
	Level string `mapstructure:"level"`
}

type listener struct {
	Prefetch        int             `mapstructure:"prefetch"`
	AcknowledgeMode AcknowledgeMode `mapstructure:"acknowledge_mode"`
}

type retry struct {
	InitialInterval time.Duration `mapstructure:"initial_interval"`
	Multiplier      int           `mapstructure:"multiplier"`
	MaxAttempts     int           `json:"max_attempts"`
}

var Conf *RabbitmqConfig

func init() {
	Conf = new(RabbitmqConfig)
	Conf.Retry = &retry{
		InitialInterval: time.Second,
		Multiplier:      2,
		MaxAttempts:     3,
	}
	Conf.Listener = &listener{
		Prefetch:        1,
		AcknowledgeMode: AcknowledgeModeAuto,
	}
	Conf.Log = &log{
		Level: "info",
	}
}
