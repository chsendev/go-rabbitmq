package config

type RabbitmqConfig struct {
	Url       string     `mapstructure:"url"`
	Log       *log       `mapstructure:"log"`
	Publisher *publisher `mapstructure:"publisher"`
	Listener  *listener  `mapstructure:"listener"`
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

type publisher struct {
	Retry        *retry   `mapstructure:"retry"`
	Confirm      *confirm `mapstructure:"confirm"`
	NotifyReturn bool     `mapstructure:"notify_return"`
}

type retry struct {
	Enabled bool `mapstructure:"enabled"`
	// InitialInterval 单位秒
	InitialInterval int `mapstructure:"initial_interval"`
	Multiplier      int `mapstructure:"multiplier"`
	MaxAttempts     int `json:"max_attempts"`
}

type confirm struct {
	Enabled   bool `mapstructure:"enabled"`
	WaitMilli int  `mapstructure:"wait_milli"`
}

var Conf *RabbitmqConfig

func Init(c *RabbitmqConfig) {
	Conf = c
}
