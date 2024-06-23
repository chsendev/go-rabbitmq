package mq

import "github.com/cscoder0/go-rabbitmq/config"

type acknowledge interface {
	Ack() error
	Nack() error
}

type acknowledgeNone struct {
}

func (a *acknowledgeNone) Ack() error {
	return nil
}
func (a *acknowledgeNone) Nack() error {
	return nil
}

type acknowledgeAuto struct {
	*rawMessage
}

func (a *acknowledgeAuto) Ack() error {
	return a.rawMessage.Ack(false)
}
func (a *acknowledgeAuto) Nack() error {
	return a.rawMessage.Nack(false, true)
}

type acknowledgeManual struct {
}

func (a *acknowledgeManual) Ack() error {
	return nil
}
func (a *acknowledgeManual) Nack() error {
	return nil
}

func isAutoAck() bool {
	if config.Conf.Listener == nil ||
		config.Conf.Listener.AcknowledgeMode == "" ||
		config.Conf.Listener.AcknowledgeMode == config.AcknowledgeModeNone {
		return true
	}
	return false
}

func getAcknowledge(msg *rawMessage) acknowledge {
	if config.Conf.Listener == nil || config.Conf.Listener.AcknowledgeMode == "" || config.Conf.Listener.AcknowledgeMode == config.AcknowledgeModeNone {
		return &acknowledgeNone{}
	} else if config.Conf.Listener.AcknowledgeMode == config.AcknowledgeModeAuto {
		return &acknowledgeAuto{rawMessage: msg}
	} else if config.Conf.Listener.AcknowledgeMode == config.AcknowledgeModeManual {
		return &acknowledgeManual{}
	} else {
		return &acknowledgeNone{}
	}
}
