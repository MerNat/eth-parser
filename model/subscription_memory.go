package model

import "sync"

type MemorySubscriptionStorage struct {
    data  map[string]bool
    mutex sync.Mutex
}

func NewMemorySubscriptionStorage() SubscriptionStorage {
    return &MemorySubscriptionStorage{data: make(map[string]bool)}
}

func (s *MemorySubscriptionStorage) Subscribe(address string) bool {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    if _, exists := s.data[address]; exists {
        return false
    }
    s.data[address] = true
    return true
}

func (s *MemorySubscriptionStorage) IsSubscribed(address string) bool {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    return s.data[address]
}
