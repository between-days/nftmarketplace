package listings

import "nftmarketplace/models"

type MockListingsAPI struct {
	// config
	ListingsHandler
	MockGetListings   func(seller string, limit int) (*[]models.Listing, bool, error)
	MockCreateListing func(newListing *models.Listing) (*models.Listing, error)
	MockDeleteListing func(listingID uint) (*models.Listing, error)
}

func (listings MockListingsAPI) GetListings(seller string, limit int) (*[]models.Listing, bool, error) {
	return listings.MockGetListings(seller, limit)
}

func (listings MockListingsAPI) CreateListing(newListing *models.Listing) (*models.Listing, error) {
	return listings.MockCreateListing(newListing)
}

func (listings MockListingsAPI) DeleteListing(listingID uint) (*models.Listing, error) {
	return listings.MockDeleteListing(listingID)
}
