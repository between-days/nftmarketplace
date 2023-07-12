package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"nftmarketplace/api"
	"nftmarketplace/db/listingsrepo"
	"nftmarketplace/listings"
	"nftmarketplace/models"
	"nftmarketplace/transactions"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const REDIS_TRANSACTION_KEY = "lastTransactionID"

type DemoConfig struct {
	webPort string
	api     *gin.Engine
	db      *gorm.DB
	cache   *redis.Client
}

func demo(demoConfig DemoConfig) {
	fmt.Print("\n\n -- DEMO MODE-- \n\n")

	// Migrate the schema
	err := demoConfig.db.AutoMigrate(&models.Listing{}, &models.Transaction{})
	if err != nil {
		panic(err)
	}

	// create some dummy data
	for i := 0; i < 60; i++ {
		ts := demoConfig.db.Create(&models.Listing{
			Model:    gorm.Model{},
			Seller:   fmt.Sprintf("fakename%d", rand.Intn(100)),
			Price:    rand.Float64() * 100,
			ImageURL: fmt.Sprintf("fakeurl%d", rand.Intn(100)),
		})
		if ts.Error != nil {
			panic(err)
		}
	}

	var ctx = context.Background()
	err = demoConfig.cache.Set(ctx, REDIS_TRANSACTION_KEY, 0, 0).Err()
	if err != nil {
		panic(err)
	}

	_, err = demoConfig.cache.Get(ctx, REDIS_TRANSACTION_KEY).Result()
	if err != nil {
		panic(err)
	}

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", demoConfig.webPort), demoConfig.api))
}

func main() {
	webPort := os.Getenv("PORT")
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	postgresDSN := os.Getenv("POSTGRES_DSN")

	db, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	cache := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	router := gin.Default()

	api := api.NewMarketPlaceAPI(
		router,
		listings.NewListingsService(cache,
			listingsrepo.NewListingsRepo(db),
		),
		transactions.NewTransactionsService(cache),
	)

	demoConfig := DemoConfig{
		webPort,
		api,
		db,
		cache,
	}

	demo(demoConfig)
}
