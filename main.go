package main

import (
	"ether-parser/api"
	"ether-parser/model"
	"ether-parser/service/parser"
)

func main() {
	subStorage := model.NewMemorySubscriptionStorage()
	txStorage := model.NewMemoryTransactionStorage()
	ethereumRPCURL := "https://ethereum-rpc.publicnode.com"
	parserService := parser.NewParser(subStorage, txStorage, ethereumRPCURL)

	// Poll blocks every 10 seconds to get the latest transactions
	parserService.PollBlocks(10)


	handler := api.NewHandler(parserService)
	api.StartServer(handler, "8080")
}
