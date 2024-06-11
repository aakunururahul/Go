package main

import "sync"

// CappedQueue stores items in FIFO order,
// but it has a capacity and will delete older items if capacity is reached.
type CappedQueue[T any] struct {
	items []T
	lock  *sync.RWMutex
}

func NewCappedQueue[T any]() *CappedQueue[T] {
	return &CappedQueue[T]{
		items: make([]T, 0),
		lock:  new(sync.RWMutex),
	}
}

func (q *CappedQueue[T]) Append(item T) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.items = append(q.items, item)
}

func (q *CappedQueue[T]) Copy() []T {
	q.lock.RLock()
	defer q.lock.RUnlock()

	copied := make([]T, len(q.items))
	copy(copied, q.items)

	return copied
}
