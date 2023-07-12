package transactionsrepo

// TODO: make repos for listings and transactions

import (
	"nftmarketplace/models"

	"gorm.io/gorm"
)

type TransactionsRepoHandler interface {
	NewTransaction(newTransaction *models.Transaction) (*models.Transaction, error)
	GetTransactions(seller, buyer string, limit int) (*[]models.Transaction, error)
}

type TransactionsRepoAPI struct {
	TransactionsRepoHandler
	PostgresClient *gorm.DB
}

func (tr TransactionsRepoAPI) NewTransaction(newTransaction *models.Transaction) (*models.Transaction, error) {
	ts := tr.PostgresClient.Create(newTransaction)
	if ts.Error != nil {
		return nil, ts.Error
	}

	return newTransaction, nil
}

func (tr TransactionsRepoAPI) GetTransactions(seller string, limit int) (*[]models.Transaction, error) {
	var ts []models.Transaction
	base := tr.PostgresClient.Where("id <> ?", -1)

	if len(seller) > 0 {
		base = base.Where("seller = ?", seller)
	}

	dbts := base.Find(&ts).Limit(limit)
	if dbts.Error != nil {
		return nil, dbts.Error
	}

	return &ts, nil
}

func NewTransactionsRepo(postgresClient *gorm.DB) TransactionsRepoAPI {
	transactionsRepo := TransactionsRepoAPI{
		PostgresClient: postgresClient,
	}

	return transactionsRepo
}
