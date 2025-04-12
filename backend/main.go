package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type BidIncomingItem struct {
	ItemId    uint32 `json:"itemId"`
	BidAmount uint32 `json:"bidAmount"`
	Bidder    uint16 `json:"bidder"`
}

type BidItem struct {
	Id        uint32 `json:"id"`
	ItemId    uint32 `json:"itemId"`
	Bidder    uint16 `json:"bidder"`
	Amount    uint32 `json:"amount"`
	TimeStamp string `json:"timeStamp"`
}

type AuctionItem struct {
	Id          uint32    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartingBid uint32    `json:"startingBid"`
	CurrentBid  uint32    `json:"currentBid"`
	Bids        []BidItem `json:"bids"`
}

var auctionItems = []AuctionItem{
	{Id: 1,
		Name:        "Vintage Watch",
		Description: "Classic timepiece from 1950s",
		StartingBid: 100,
		CurrentBid:  100,
		Bids:        []BidItem{},
	},
	{Id: 2,
		Name:        "Art Painting",
		Description: "Original Artwork in Skribble.io",
		StartingBid: 150,
		CurrentBid:  150,
		Bids:        []BidItem{},
	},
}

var bidHistory []BidItem

// greeting for default page.
func greeting(c *gin.Context) {
	c.String(http.StatusOK, "MLH GHW API Week!")
}

func fetchAuctionItems(c *gin.Context) {
	c.JSON(http.StatusOK, auctionItems)
}

func fetchAuctionItemById(c *gin.Context) {
	idParam := c.Param("id")

	// Convert string ID to int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var foundAuctionItem *AuctionItem

	for _, auctionItem := range auctionItems {
		if int(id) == int(auctionItem.Id) {
			foundAuctionItem = &auctionItem
			break
		}
	}

	if foundAuctionItem != nil {
		c.JSON(http.StatusOK, foundAuctionItem)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
	}
}

func processAuctionBids(c *gin.Context) {
	var newBidItem BidIncomingItem

	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields() // Reject unknown fields

	// Call Decode to bind the received JSON to
	// newBidItem.
	if err := decoder.Decode(&newBidItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if newBidItem.ItemId > 0 && newBidItem.BidAmount > 0 && newBidItem.Bidder > 0 {
		// valid input data
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data!"})
		return
	}

	id := newBidItem.ItemId

	var foundAuctionItem *AuctionItem

	/*
		for _, auctionItem := range auctionItems {
			if int(id) == int(auctionItem.Id) {
				foundAuctionItem = &auctionItem
				break
			}
		}
	*/
	for i := range auctionItems {
		if int(id) == int(auctionItems[i].Id) {
			foundAuctionItem = &auctionItems[i]
			break
		}
	}

	if foundAuctionItem == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	if newBidItem.BidAmount <= foundAuctionItem.CurrentBid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bid must be higher than current bid"})
		return
	}

	/*
		for i := range auctionItems {
			if int(id) == int(auctionItems[i].Id) {
				auctionItems[i].CurrentBid = newBidItem.BidAmount
				break
			}
		}
	*/

	// Update CurrentBid with new bid amount
	foundAuctionItem.CurrentBid = newBidItem.BidAmount

	newId := len(bidHistory) + 1
	currDate := time.Now().UTC().Format(time.RFC3339)
	bidItem := BidItem{
		Id:        uint32(newId),
		ItemId:    newBidItem.ItemId,
		Bidder:    newBidItem.Bidder,
		Amount:    newBidItem.BidAmount,
		TimeStamp: currDate,
	}

	bidHistory = append(bidHistory, bidItem)

	// Update Bids with new bid
	foundAuctionItem.Bids = append(foundAuctionItem.Bids, bidItem)

	log.Println("bid history:", bidHistory)

	c.JSON(http.StatusCreated, bidItem)
}

// Helper function to get and validate an environment variable
func getEnvVar(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("Error: %s environment variable is not set", key)
	}
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return "", fmt.Errorf("Error: %s is empty or contains only spaces", key)
	}
	return trimmedValue, nil
}

func getEnvironmentVariables() (string, error) {
	// Load env vars
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("Error loading .env file: %w", err)
	}

	serverAddr, err := getEnvVar("SERVER_ADDR")
	if err != nil {
		return "", err
	}

	return serverAddr, nil
}

func main() {
	// get env vars
	serverAddr, err := getEnvironmentVariables()
	if err != nil {
		log.Fatal("Failed to load environment variables:", err)
	}

	router := gin.Default()
	router.GET("/", greeting)
	router.GET("/api/items", fetchAuctionItems)
	router.GET("/api/items/:id", fetchAuctionItemById)
	router.POST("/api/bids", processAuctionBids)

	router.Run(serverAddr)
}
