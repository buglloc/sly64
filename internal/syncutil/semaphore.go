package syncutil

import "context"

type Semaphore interface {
	Acquire(ctx context.Context) (err error)
	Release()
}

var _ Semaphore = (*NopSemaphore)(nil)

type NopSemaphore struct{}

func (NopSemaphore) Acquire(_ context.Context) error {
	return nil
}

func (NopSemaphore) Release() {}

var _ Semaphore = (*LeakySemaphore)(nil)

type LeakySemaphore struct {
	c chan struct{}
}

func NewLeakySemaphore(limit int) *LeakySemaphore {
	return &LeakySemaphore{
		c: make(chan struct{}, limit),
	}
}

func (c *LeakySemaphore) Acquire(ctx context.Context) (err error) {
	select {
	case c.c <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *LeakySemaphore) Release() {
	select {
	case <-c.c:
	default:
	}
}
