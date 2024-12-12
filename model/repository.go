package model

type SubscriptionStorage interface {
	Subscribe(address string) bool
	IsSubscribed(address string) bool
}

type TransactionStorage interface {
	AddTransaction(address string, transaction Transaction)
	GetTransactions(address string) []Transaction
}
