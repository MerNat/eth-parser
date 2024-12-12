package api

import (
	"fmt"
	"log"
	"net/http"
)

func SetupRouter(handler *Handler) {
	http.HandleFunc("/subscribe", handler.Subscribe)
	http.HandleFunc("/transactions", handler.Transactions)
	http.HandleFunc("/currentBlock", handler.CurrentBlock)
}

func StartServer(handler *Handler, port string) {
	SetupRouter(handler)
	log.Printf("Server started on :%s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
