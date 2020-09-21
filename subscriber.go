package reactor

import (
	"context"
	"math"
)

const RequestInfinite = math.MaxInt32

var emptySubscriber = &subscriber{}

type Subscription interface {
	Request(n int)
	Cancel()
}

type Subscriber interface {
	OnComplete()
	OnError(error)
	OnNext(Any)
	OnSubscribe(context.Context, Subscription)
}

type subscriber struct {
	fnOnSubscribe FnOnSubscribe
	fnOnNext      FnOnNext
	fnOnComplete  FnOnComplete
	fnOnError     FnOnError
}

func (p *subscriber) OnComplete() {
	if p == nil || p.fnOnComplete == nil {
		return
	}
	p.fnOnComplete()
}

func (p *subscriber) OnError(err error) {
	if p == nil || p.fnOnError == nil {
		return
	}
	p.fnOnError(err)
}

func (p *subscriber) OnSubscribe(ctx context.Context, s Subscription) {
	if p == nil || p.fnOnSubscribe == nil {
		s.Request(RequestInfinite)
	} else {
		p.fnOnSubscribe(ctx, s)
	}
}

func (p *subscriber) OnNext(i Any) {
	if p.fnOnNext == nil {
		return
	}
	if err := p.fnOnNext(i); err != nil {
		p.OnError(err)
	}
}

type SubscriberOption func(*subscriber)

func OnNext(onNext FnOnNext) SubscriberOption {
	return func(s *subscriber) {
		s.fnOnNext = onNext
	}
}

func OnComplete(onComplete FnOnComplete) SubscriberOption {
	return func(s *subscriber) {
		s.fnOnComplete = onComplete
	}
}

func OnError(onError FnOnError) SubscriberOption {
	return func(i *subscriber) {
		i.fnOnError = onError
	}
}

func OnSubscribe(onSubscribe FnOnSubscribe) SubscriberOption {
	return func(i *subscriber) {
		i.fnOnSubscribe = onSubscribe
	}
}

func NewSubscriber(opts ...SubscriberOption) Subscriber {
	if len(opts) < 1 {
		return emptySubscriber
	}
	s := &subscriber{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
