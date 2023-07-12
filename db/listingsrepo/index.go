package listingsrepo

// TODO: make repos for listings and transactions

import (
	"nftmarketplace/models"

	"gorm.io/gorm"
)

type ListingsRepoHandler interface {
	GetListings(seller string, limit int) (*[]models.Listing, error)
	CreateListing(newListing *models.Listing) (*models.Listing, error)
	DeleteListing(listingID uint) (*models.Listing, error)
}

type ListingsRepoAPI struct {
	ListingsRepoHandler
	PostgresClient *gorm.DB
}

func (lr ListingsRepoAPI) GetListings(seller string, limit int) (*[]models.Listing, error) {
	var ls []models.Listing
	base := lr.PostgresClient.Where("id <> ?", -1)

	if len(seller) > 0 {
		base = base.Where("seller = ?", seller)
	}

	ts := base.Find(&ls).Limit(limit)
	if ts.Error != nil {
		return nil, ts.Error
	}

	return &ls, nil

}

func (lr ListingsRepoAPI) CreateListing(newListing *models.Listing) (*models.Listing, error) {
	ts := lr.PostgresClient.Create(newListing)
	if ts.Error != nil {
		// TODO: listings/transactions error structs for accurate status from calling api
		// if ts.Error == gorm.ErrRecordNotFound{
		// 	return nil,
		// }
		return nil, ts.Error
	}

	return newListing, nil
}

func (lr ListingsRepoAPI) DeleteListing(listingID uint) (*models.Listing, error) {
	ts := lr.PostgresClient.Delete(&models.Listing{}, 1)

	if ts.Error != nil {
		return nil, ts.Error
	}

	var deletedListing models.Listing
	ts = lr.PostgresClient.Unscoped().First(&deletedListing, listingID).Limit(30)
	if ts.Error != nil {
		return nil, ts.Error
	}

	return &deletedListing, nil
}

func NewListingsRepo(postgresClient *gorm.DB) ListingsRepoAPI {
	listingsRepo := ListingsRepoAPI{
		PostgresClient: postgresClient,
	}

	return listingsRepo
}
