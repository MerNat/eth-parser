package parser

import (
	"bytes"
	"encoding/json"
	"ether-parser/model"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ParserService struct {
	currentBlock   int
	subStorage     model.SubscriptionStorage
	txStorage      model.TransactionStorage
	ethereumRPCURL string
	mutex          sync.Mutex
}

func NewParser(subStorage model.SubscriptionStorage, txStorage model.TransactionStorage, ethereumRPCURL string) Parser {
	return &ParserService{
		currentBlock:   0,
		subStorage:     subStorage,
		txStorage:      txStorage,
		ethereumRPCURL: ethereumRPCURL,
	}
}

func (p *ParserService) GetCurrentBlock() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.currentBlock
}

func (p *ParserService) Subscribe(address string) bool {
	return p.subStorage.Subscribe(strings.ToLower(address))
}

func (p *ParserService) GetTransactions(address string) []model.Transaction {
	return p.txStorage.GetTransactions(strings.ToLower(address))
}

func (p *ParserService) FetchCurrentBlock() (int, error) {
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

	resp, err := http.Post(p.ethereumRPCURL, "application/json", bytes.NewReader(body))
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

func (p *ParserService) FetchBlockTransactions(blockNumber int) ([]model.Transaction, error) {
	type Response struct {
		Result struct {
			Transactions []model.Transaction `json:"transactions"`
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

	resp, err := http.Post(p.ethereumRPCURL, "application/json", bytes.NewReader(body))
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

func (p *ParserService) ProcessBlock(blockNumber int) {
	transactions, err := p.FetchBlockTransactions(blockNumber)
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

func (p *ParserService) PollBlocks(intervalSeconds int) {
	go func() {
		for {
			latestBlock, err := p.FetchCurrentBlock()
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
				p.ProcessBlock(block)
			}

			time.Sleep(time.Duration(intervalSeconds) * time.Second)
		}
	}()
}
