package transactions

import (
	"nftmarketplace/models"
)

type MockTransactionsAPI struct {
	TransactionsHandler
	MockNewTransaction func(*models.Transaction) (*models.Transaction, error)
}

func (transactions MockTransactionsAPI) NewTransaction(ts *models.Transaction) (*models.Transaction, error) {
	return transactions.MockNewTransaction(ts)
}
