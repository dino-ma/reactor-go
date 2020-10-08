package internal

import (
	"errors"
	"fmt"
)

var ErrCallOnSubscribeDuplicated = errors.New("call OnSubscribe duplicated")
var EmptySubscription = emptySubscription{}

func TryRecoverError(re interface{}) error {
	if re == nil {
		return nil
	}
	switch e := re.(type) {
	case error:
		return e
	case string:
		return errors.New(e)
	default:
		return fmt.Errorf("%s", e)
	}
}

type emptySubscription struct {
}

func (emptySubscription) Request(n int) {
}

func (emptySubscription) Cancel() {
}

func SafeCloseDone(done chan<- struct{}) (ok bool) {
	defer func() {
		ok = recover() == nil
	}()
	close(done)
	return
}
