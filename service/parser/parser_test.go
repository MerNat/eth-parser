package parser

import (
	"ether-parser/model"
	"testing"
)

// Testing the parser service using mocking to simulate the real logic implementation
type MockParser struct {
	currentBlock int
	subStorage   *model.MemorySubscriptionStorage
	txStorage    *model.MemoryTransactionStorage
}

func NewMockParser() *MockParser {
	return &MockParser{
		currentBlock: 0,
		subStorage:   model.NewMemorySubscriptionStorage().(*model.MemorySubscriptionStorage),
		txStorage:    model.NewMemoryTransactionStorage().(*model.MemoryTransactionStorage),
	}
}

func (m *MockParser) FetchCurrentBlock() (int, error) {
	return 100, nil // Simulate a block number
}

func (m *MockParser) FetchBlockTransactions(blockNumber int) ([]model.Transaction, error) {
	return []model.Transaction{
		{Hash: ptr("0xabc"), From: ptr("0x123"), To: ptr("0x456"), Value: ptr("100"), BlockNumber: ptr("100")},
	}, nil // Simulate transactions
}

func (m *MockParser) ProcessBlock(blockNumber int) {
	transactions, _ := m.FetchBlockTransactions(blockNumber)
	for _, tx := range transactions {
		m.txStorage.AddTransaction(*tx.From, tx)
	}
	m.currentBlock = blockNumber
}

func TestParserService(t *testing.T) {
	mockParser := NewMockParser()

	// Test subscribing
	if !mockParser.subStorage.Subscribe("0x123") {
		t.Errorf("Expected subscription to succeed")
	}

	// Test fetching current block
	block, err := mockParser.FetchCurrentBlock()
	if err != nil || block != 100 {
		t.Errorf("Expected block to be 100, got %d, err: %v", block, err)
	}

	// Test processing block
	mockParser.ProcessBlock(100)
	transactions := mockParser.txStorage.GetTransactions("0x123")
	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(transactions))
	}
	if *transactions[0].Hash != "0xabc" {
		t.Errorf("Expected transaction hash to be '0xabc', got '%s'", *transactions[0].Hash)
	}
}

func ptr(value string) *string {
	return &value
}
