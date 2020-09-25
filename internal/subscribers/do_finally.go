package subscribers

import (
	"context"
	"sync"

	"github.com/jjeffcaii/reactor-go"
)

type DoFinallySubscriber struct {
	actual    reactor.Subscriber
	onFinally reactor.FnOnFinally
	once      sync.Once
	s         reactor.Subscription
}

func (p *DoFinallySubscriber) Request(n int) {
	p.s.Request(n)
}

func (p *DoFinallySubscriber) Cancel() {
	p.s.Cancel()
	p.runFinally(reactor.SignalTypeCancel)
}

func (p *DoFinallySubscriber) OnError(err error) {
	p.actual.OnError(err)
	if reactor.IsCancelledError(err) {
		p.runFinally(reactor.SignalTypeCancel)
	} else {
		p.runFinally(reactor.SignalTypeError)
	}
}

func (p *DoFinallySubscriber) OnNext(v reactor.Any) {
	p.actual.OnNext(v)
}

func (p *DoFinallySubscriber) OnSubscribe(ctx context.Context, s reactor.Subscription) {
	p.s = s
	p.actual.OnSubscribe(ctx, p)
}

func (p *DoFinallySubscriber) OnComplete() {
	p.actual.OnComplete()
	p.runFinally(reactor.SignalTypeComplete)
}

func (p *DoFinallySubscriber) runFinally(sig reactor.SignalType) {
	p.once.Do(func() {
		p.onFinally(sig)
	})
}

func NewDoFinallySubscriber(actual reactor.Subscriber, onFinally reactor.FnOnFinally) *DoFinallySubscriber {
	return &DoFinallySubscriber{
		onFinally: onFinally,
		actual:    actual,
	}
}
