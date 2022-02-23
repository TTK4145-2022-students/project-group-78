package events

import "sync"

type Topic struct {
	subscribers []chan Event
	mutex       sync.Mutex
}

func (s *Topic) Publish(event Event) {
	s.mutex.Lock()
	for _, subscriber := range s.subscribers {
		subscriber <- event
	}
	s.mutex.Unlock()
}

func (s *Topic) Subscribe(subscriber chan Event) {
	s.mutex.Lock()
	s.subscribers = append(s.subscribers, subscriber)
	s.mutex.Unlock()
}

func (s *Topic) Unsubscribe(subscriber chan Event) {
	s.mutex.Lock()
	for i := 0; i < len(s.subscribers); i++ {
		if s.subscribers[i] == subscriber {
			s.subscribers = append(s.subscribers[:i], s.subscribers[i+1:]...) // Idiomatic/idiotic way of deleting element in slice
		}
	}
	s.mutex.Unlock()
}
