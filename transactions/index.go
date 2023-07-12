package transactions

import (
	"context"
	"fmt"
	"nftmarketplace/models"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const REDIS_TRANSACTION_KEY = "lastTransactionID"

type TransactionsHandler interface {
	NewTransaction(ts *models.Transaction) (*models.Transaction, error)
}

type TransactionsAPI struct {
	RedisClient *redis.Client
	TransactionsHandler
}

func (transactions TransactionsAPI) NewTransaction(newTransaction *models.Transaction) (*models.Transaction, error) {

	// pull relevant listing

	// check buyer has enough in wallet

	var ctx = context.Background()

	// TODO: fix cachine
	// caching not concurrent safe, nothing stopping a request getting increment during new transaction
	// should use redis increment command ~transactions.RedisClient.Incr()

	id, err := transactions.RedisClient.Get(ctx, REDIS_TRANSACTION_KEY).Result()
	if err != nil {
		return nil, err
	}

	ID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}

	ID += 1
	newTransaction.ID = uint(ID)

	err = transactions.RedisClient.Set(ctx, REDIS_TRANSACTION_KEY, fmt.Sprint(ID), 0).Err()
	if err != nil {
		return nil, err
	}

	// TODO:
	// to avoid duplicate ids on transactions, make sure to fail after increment updated
	// external service will pick up incongruency

	return newTransaction, nil
}

func NewTransactionsService(redisClient *redis.Client) TransactionsAPI {
	var ts TransactionsAPI
	ts.RedisClient = redisClient
	return ts
}
