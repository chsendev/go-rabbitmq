package publish

import (
	"time"
)

func WithConfirm(wait time.Duration) Opt {
	return func(p *Publisher) {
		p.confirm = true
		p.waitMilli = wait
	}
}

func WithNotifyReturn() Opt {
	return func(p *Publisher) {
		p.Client.Channel.NotifyReturn(notifyReturn)
		p.mandatory = true
	}
}
func WithDelay(delayed time.Duration) Opt {
	return func(p *Publisher) {
		p.headers["x-delay"] = int(delayed / time.Millisecond)
	}
}
