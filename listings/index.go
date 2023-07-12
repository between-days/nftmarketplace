package listings

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"time"

	"nftmarketplace/db/listingsrepo"
	"nftmarketplace/models"

	"github.com/go-redis/redis/v8"
)

// TODO: move db to own layer as repo pattern

type ListingsHandler interface {
	GetListings(seller string, limit int) (*[]models.Listing, bool, error)
	CreateListing(newListing *models.Listing) (*models.Listing, error)
	DeleteListing(listingID uint) (*models.Listing, error)
}

type ListingsAPI struct {
	ListingsHandler
	RedisClient *redis.Client

	ListingsRepo listingsrepo.ListingsRepoHandler
}

func (listings ListingsAPI) GetListings(seller string, limit int) (*[]models.Listing, bool, error) {
	var ls *[]models.Listing
	cache := false

	qs := fmt.Sprintf("seller=%s", seller)

	rls, err := listings.RedisClient.Get(context.Background(), qs).Result()
	if err != nil {
		// TODO: paginate, add options
		// TODO: <> probably not good for performance. fix later

		ls, err = listings.ListingsRepo.GetListings(seller, limit)
		if err != nil {
			return nil, cache, err
		}

		lsb, err := json.Marshal(&ls)
		if err != nil {
			log.Printf("failed marshalling db response for storage in cache: %s\n", err)
		}

		err = listings.RedisClient.Set(context.Background(), qs, lsb, time.Second*10).Err()
		if err != nil {
			log.Printf("failed storing db result in cache: %s\n", err)
		} else {
			log.Printf("refreshed cache on get listings ( seller=%s )\n", seller)
		}

	} else {
		err := json.Unmarshal([]byte(rls), &ls)
		if err != nil {
			return nil, true, err
		}
		cache = true
	}

	return ls, cache, nil
}

func (listings ListingsAPI) CreateListing(newListing *models.Listing) (*models.Listing, error) {
	// TODO: only owner on blockchain can create listing
	ls, err := listings.ListingsRepo.CreateListing(newListing)
	if err != nil {
		return nil, err
	}

	return ls, nil
}

// TODO: update listing
// only owner or admin can update listing

// TODO: changes to listing need to be recorded incase of tampering
// logging service or chuck in strictly writable (no updates to records) db

func (listings ListingsAPI) DeleteListing(listingID uint) (*models.Listing, error) {
	// TODO: only seller or admin can delete listing

	ls, err := listings.ListingsRepo.DeleteListing(listingID)

	if err != nil {
		return nil, err
	}

	return ls, nil
}

func NewListingsService(redisClient *redis.Client, listingsRepo listingsrepo.ListingsRepoHandler) ListingsAPI {
	var ls ListingsAPI
	ls.RedisClient = redisClient
	// ls.PostgresClient = postgresClient
	ls.ListingsRepo = listingsRepo
	return ls
}
