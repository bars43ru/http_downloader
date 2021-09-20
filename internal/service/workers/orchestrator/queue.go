package orchestrator

import (
	"container/list"
	"sync"
)

type queue struct {
	lock     sync.Mutex
	items    list.List
}

func (q *queue) Enqueue(items ...string) {
	q.lock.Lock()
	defer q.lock.Unlock()
	for _, v := range items {
		q.items.PushBack(v)
	}
}

func (q *queue) TryDequeue() (string, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.items.Len() > 0 {
		v := q.items.Front()
		q.items.Remove(v)
		return v.Value.(string), true
	}
	return "", false
}