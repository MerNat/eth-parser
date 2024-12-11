package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const EthereumRPCURL = "https://ethereum-rpc.publicnode.com"

type Transaction struct {
	Hash        *string `json:"hash"`
	From        *string `json:"from"`
	To          *string `json:"to"`
	Value       *string `json:"value"`
	BlockNumber *string `json:"blockNumber"`
}

type SubscriptionStorage interface {
	Subscribe(address string) bool
	IsSubscribed(address string) bool
}

type TransactionStorage interface {
	AddTransaction(address string, transaction Transaction)
	GetTransactions(address string) []Transaction
}

type MemorySubscriptionStorage struct {
	data  map[string]bool
	mutex sync.Mutex
}

func NewMemorySubscriptionStorage() *MemorySubscriptionStorage {
	return &MemorySubscriptionStorage{
		data: make(map[string]bool),
	}
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

type MemoryTransactionStorage struct {
	data  map[string][]Transaction
	mutex sync.Mutex
}

func NewMemoryTransactionStorage() *MemoryTransactionStorage {
	return &MemoryTransactionStorage{
		data: make(map[string][]Transaction),
	}
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

type Parser struct {
	currentBlock int
	subStorage   SubscriptionStorage
	txStorage    TransactionStorage
	mutex        sync.Mutex
}

func NewParser() *Parser {
	return &Parser{
		currentBlock: 0,
		subStorage:   NewMemorySubscriptionStorage(),
		txStorage:    NewMemoryTransactionStorage(),
	}
}

func (p *Parser) GetCurrentBlock() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.currentBlock
}

func (p *Parser) Subscribe(address string) bool {
	return p.subStorage.Subscribe(strings.ToLower(address))
}

func (p *Parser) GetTransactions(address string) []Transaction {
	return p.txStorage.GetTransactions(strings.ToLower(address))
}

func (p *Parser) fetchCurrentBlock() (int, error) {
	type Response struct {
		Result string `json:"result"`
	}

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      1,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(EthereumRPCURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response Response
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return 0, err
	}

	blockNumber, err := strconv.ParseInt(response.Result, 0, 64)
	if err != nil {
		return 0, err
	}

	return int(blockNumber), nil
}

func (p *Parser) fetchBlockTransactions(blockNumber int) ([]Transaction, error) {
	type Response struct {
		Result struct {
			Transactions []Transaction `json:"transactions"`
		} `json:"result"`
	}

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{fmt.Sprintf("0x%x", blockNumber), true},
		"id":      1,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(EthereumRPCURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response Response
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, err
	}

	return response.Result.Transactions, nil
}

func (p *Parser) processBlock(blockNumber int) {
	transactions, err := p.fetchBlockTransactions(blockNumber)
	if err != nil {
		fmt.Println("Error fetching transactions for block:", blockNumber, "Error:", err)
		return
	}

	for _, tx := range transactions {
		if tx.From != nil {
			fromAddress := strings.ToLower(*tx.From)
			if p.subStorage.IsSubscribed(fromAddress) {
				p.txStorage.AddTransaction(fromAddress, tx)
			}
		}
		if tx.To != nil {
			toAddress := strings.ToLower(*tx.To)
			if p.subStorage.IsSubscribed(toAddress) {
				p.txStorage.AddTransaction(toAddress, tx)
			}
		}
	}

	p.mutex.Lock()
	p.currentBlock = blockNumber
	p.mutex.Unlock()
}

func (p *Parser) PollBlocks(intervalSeconds int) {
	go func() {
		for {
			latestBlock, err := p.fetchCurrentBlock()
			if err != nil {
				fmt.Println("Error fetching latest block:", err)
				time.Sleep(time.Duration(intervalSeconds) * time.Second)
				continue
			}

			p.mutex.Lock()
			currentBlock := p.currentBlock
			if currentBlock == 0 {
				p.currentBlock = latestBlock // Initialize currentBlock if running for the first time
				log.Printf("Initialized current block to: %d\n", latestBlock)
				p.mutex.Unlock()
				time.Sleep(time.Duration(intervalSeconds) * time.Second)
				continue
			}
			p.mutex.Unlock()

			// Process new blocks
			for block := currentBlock + 1; block <= latestBlock; block++ {
				log.Printf("Processing new block: %d\n", block)
				p.processBlock(block)
			}

			time.Sleep(time.Duration(intervalSeconds) * time.Second)
		}
	}()
}

func main() {
	parser := NewParser()

	fmt.Println("Starting Ethereum Transaction Parser...")

	// Poll for new blocks every 10 seconds
	parser.PollBlocks(10)

	// API Endpoints
	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Address string `json:"address"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		if parser.Subscribe(req.Address) {
			latestBlock, err := parser.fetchCurrentBlock()
			if err != nil {
				http.Error(w, "Error fetching latest block", http.StatusInternalServerError)
				return
			}
			parser.processBlock(latestBlock)
			fmt.Fprintf(w, "Subscribed to: %s", req.Address)
		} else {
			fmt.Fprintf(w, "Address already subscribed: %s", req.Address)
		}
	})

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		address := strings.ToLower(r.URL.Query().Get("address"))
		if address == "" {
			http.Error(w, "Address query parameter is required", http.StatusBadRequest)
			return
		}
		transactions := parser.GetTransactions(address)
		if err := json.NewEncoder(w).Encode(transactions); err != nil {
			http.Error(w, "Error encoding transactions", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/currentBlock", func(w http.ResponseWriter, r *http.Request) {
		currentBlock := parser.GetCurrentBlock()
		fmt.Fprintf(w, "Current Block: %d", currentBlock)
	})

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
