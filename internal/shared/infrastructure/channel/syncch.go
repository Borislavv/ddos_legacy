package channel

import "sync"

type SyncCh struct {
	mu     *sync.Mutex
	ch     chan interface{}
	isOpen bool
}

func NewSyncCh(ch chan interface{}) *SyncCh {
	return &SyncCh{
		mu:     &sync.Mutex{},
		ch:     ch,
		isOpen: true,
	}
}

func (s *SyncCh) Get() chan interface{} {
	return s.ch
}

func (s *SyncCh) Write(data string) bool {
	defer s.mu.Unlock()
	s.mu.Lock()

	if s.isOpen {
		s.ch <- data
	}
	return s.isOpen
}

func (s *SyncCh) IsOpen() bool {
	defer s.mu.Unlock()
	s.mu.Lock()
	return s.isOpen
}

func (s *SyncCh) Close() bool {
	defer s.mu.Unlock()
	s.mu.Lock()

	if s.isOpen {
		s.isOpen = false
		close(s.ch)
	}
	return !s.isOpen
}
