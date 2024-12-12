package model

import "sync"

type Transaction struct {
	Hash        *string `json:"hash"`
	From        *string `json:"from"`
	To          *string `json:"to"`
	Value       *string `json:"value"`
	BlockNumber *string `json:"blockNumber"`
}

type MemoryTransactionStorage struct {
	data  map[string][]Transaction
	mutex sync.Mutex
}

func NewMemoryTransactionStorage() TransactionStorage {
	return &MemoryTransactionStorage{data: make(map[string][]Transaction)}
}

func (s *MemoryTransactionStorage) AddTransaction(address string, transaction Transaction) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[address] = append(s.data[address], transaction)
}

func (s *MemoryTransactionStorage) GetTransactions(address string) []Transaction {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.data[address]
}
