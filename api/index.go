package api

import (
	"net/http"
	"nftmarketplace/listings"
	"nftmarketplace/transactions"

	"github.com/gin-gonic/gin"
)

type MarketPlaceHandler struct {
	Listings     listings.ListingsHandler
	Transactions transactions.TransactionsHandler
}

func (mph MarketPlaceHandler) hello(c *gin.Context) {
	c.String(http.StatusOK, "hello from marketplace api")
}

func NewMarketPlaceAPI(r *gin.Engine, listingsHandler listings.ListingsHandler, transactionsHandler transactions.TransactionsHandler) *gin.Engine {
	var mph MarketPlaceHandler
	mph.Listings = listingsHandler
	mph.Transactions = transactionsHandler

	v1 := r.Group("v1")
	{
		v1.GET("/hello", mph.hello)
		v1.GET("/listings/", mph.getListings)
		v1.POST("/listings/", mph.createListing)
		v1.DELETE("/listings/:id", mph.deleteListing)

		v1.POST("/transactions", mph.newTransaction)
	}
	return r
}
