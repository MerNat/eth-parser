package api

import (
	"encoding/json"
	"ether-parser/service/parser"
	"net/http"
	"strings"
)

type Handler struct {
	Parser parser.Parser
}

func NewHandler(parser parser.Parser) *Handler {
	return &Handler{Parser: parser}
}

func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	subscribed := h.Parser.Subscribe(req.Address)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]bool{"subscribed": subscribed}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *Handler) Transactions(w http.ResponseWriter, r *http.Request) {
	address := strings.ToLower(r.URL.Query().Get("address"))
	if address == "" {
		http.Error(w, "Address query parameter is required", http.StatusBadRequest)
		return
	}
	transactions := h.Parser.GetTransactions(address)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "Error encoding transactions", http.StatusInternalServerError)
	}
}

func (h *Handler) CurrentBlock(w http.ResponseWriter, r *http.Request) {
	currentBlock := h.Parser.GetCurrentBlock()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]int{"currentBlock": currentBlock}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding current block", http.StatusInternalServerError)
	}
}
