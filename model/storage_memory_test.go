package model

import (
	"testing"
)

func TestMemorySubscriptionStorage(t *testing.T) {
	storage := NewMemorySubscriptionStorage()

	// Test subscribing a new address
	if !storage.Subscribe("0x123") {
		t.Errorf("Expected subscription to succeed for new address")
	}

	// Test subscribing an existing address
	if storage.Subscribe("0x123") {
		t.Errorf("Expected subscription to fail for existing address")
	}

	// Test checking if an address is subscribed
	if !storage.IsSubscribed("0x123") {
		t.Errorf("Expected address to be subscribed")
	}
}

func TestMemoryTransactionStorage(t *testing.T) {
	storage := NewMemoryTransactionStorage()
	tx := Transaction{
		Hash:        ptr("0xabc"),
		From:        ptr("0x123"),
		To:          ptr("0x456"),
		Value:       ptr("100"),
		BlockNumber: ptr("10"),
	}

	// Test adding a transaction
	storage.AddTransaction("0x123", tx)
	transactions := storage.GetTransactions("0x123")
	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(transactions))
	}

	// Test transaction content
	if *transactions[0].Hash != "0xabc" {
		t.Errorf("Expected transaction hash to be '0xabc', got '%s'", *transactions[0].Hash)
	}
}

func ptr(value string) *string {
	return &value
}
