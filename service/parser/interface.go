package parser

import "ether-parser/model"

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []model.Transaction
	ProcessBlock(blockNumber int)
    FetchBlockTransactions(blockNumber int) ([]model.Transaction, error)
    FetchCurrentBlock() (int, error)
    PollBlocks(intervalSeconds int)
}
