package orchestrator

import (
	"context"
	"time"
)

type PubSub struct {
	ctx   context.Context
	queue queue
}

func New(ctx context.Context) *PubSub {
	return &PubSub{
		ctx:   ctx,
		queue: queue{},
	}
}

func (q *PubSub) Pub(items ...string) error {
	select {
	case <-q.ctx.Done():
		return ErrStopped
	default:
		q.queue.Enqueue(items...)
		return nil
	}
}

func (q *PubSub) Sub(out chan<- string) {
	for {
		v, ok := q.queue.TryDequeue()
		if ok {
			out <- v
		} else {
			select {
			case <-q.ctx.Done():
				close(out)
				return
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}
