package ctxutil

import (
	"context"
	"fmt"
	"time"
)

func Detach(parent context.Context) context.Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	return detachedCtx{parent: parent}
}

var _ context.Context = detachedCtx{}

type detachedCtx struct {
	parent context.Context
}

func (c detachedCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c detachedCtx) Done() <-chan struct{} {
	return nil
}

func (c detachedCtx) Err() error {
	return nil
}

func (c detachedCtx) Value(key interface{}) interface{} {
	return c.parent.Value(key)
}

func (c detachedCtx) String() string {
	return fmt.Sprintf("%s.Detach", c.parent)
}
