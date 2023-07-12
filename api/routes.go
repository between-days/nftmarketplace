package api

import (
	"fmt"
	"net/http"
	"nftmarketplace/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

const XCACHE_HEADER_KEY = "X-Cache"
const CACHE_HIT = "HIT"
const CACHE_MISS = "MISS"

func (mph MarketPlaceHandler) getListings(c *gin.Context) {
	seller := c.Query("seller")
	ls, hit, err := mph.Listings.GetListings(seller, 30)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	if hit {
		c.Header(XCACHE_HEADER_KEY, CACHE_HIT)
	} else {
		c.Header(XCACHE_HEADER_KEY, CACHE_MISS)
	}

	c.JSON(http.StatusOK, gin.H{
		"listings": ls,
	})
}

func (mph MarketPlaceHandler) createListing(c *gin.Context) {
	var newListing models.Listing
	if err := c.ShouldBindJSON(&newListing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ls, err := mph.Listings.CreateListing(&newListing)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"listing": ls,
	})
}

func (mph MarketPlaceHandler) deleteListing(c *gin.Context) {
	id := c.Param("id")
	ID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = mph.Listings.DeleteListing(uint(ID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (mph MarketPlaceHandler) newTransaction(c *gin.Context) {
	var newTransaction models.Transaction
	if err := c.ShouldBindJSON(&newTransaction); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ts, err := mph.Transactions.NewTransaction(&newTransaction)
	if err != nil {
		fmt.Printf("error: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction": ts,
	})
}
