package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type BidIncomingItem struct {
	ItemId    uint32 `json:"itemId"`
	BidAmount uint32 `json:"bidAmount"`
	Bidder    string `json:"bidder"`
}

type BidItem struct {
	Id        uint32 `json:"id"`
	ItemId    uint32 `json:"itemId"`
	Bidder    string `json:"bidder"`
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

type WSMessage struct {
	Type      string `json:"type"`
	ItemId    uint32 `json:"itemId,omitempty"`
	BidAmount uint32 `json:"bidAmount,omitempty"`
	Bidder    string `json:"bidder,omitempty"`
}

var clients = make(map[*websocket.Conn]bool)
var clientsMutex sync.Mutex
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins (use with caution in production)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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

/**
 * ======================== REST API ROUTES ========================
 * These routes demonstrate HTTP request/response communication
 */

// index for default page.
func index(c *gin.Context) {
	c.String(http.StatusOK, "MLH GHW API Week!")
}

// fetch auction items
func fetchAuctionItems(c *gin.Context) {
	// respond with existing auction items
	c.JSON(http.StatusOK, auctionItems)
}

// fetch auction item by id
func fetchAuctionItemById(c *gin.Context) {
	// get the id parameter
	idParam := c.Param("id")

	// Convert string ID to int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	foundAuctionItem := searchForAuctionItemById(uint32(id))

	if foundAuctionItem != nil {
		c.JSON(http.StatusOK, foundAuctionItem)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
	}
}

func fetchBidHistory(c *gin.Context) {
	// respond with existing bid history
	if len(bidHistory) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No bid history"})
	} else {
		c.JSON(http.StatusOK, bidHistory)
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

	// Validate input data
	if newBidItem.ItemId == 0 || newBidItem.BidAmount == 0 || len(newBidItem.Bidder) == 0 {
		// invalid input data
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data!"})
		return
	}

	id := newBidItem.ItemId

	foundAuctionItem := searchForAuctionItemById(id)

	if foundAuctionItem == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Validate the amounts
	if newBidItem.BidAmount <= foundAuctionItem.CurrentBid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bid must be higher than current bid"})
		return
	}

	// Update CurrentBid with new bid amount
	foundAuctionItem.CurrentBid = newBidItem.BidAmount

	newId := len(bidHistory) + 1
	currDate := time.Now().UTC().Format(time.RFC3339)

	// Create new bid
	newBid := BidItem{
		Id:        uint32(newId),
		ItemId:    newBidItem.ItemId,
		Bidder:    newBidItem.Bidder,
		Amount:    newBidItem.BidAmount,
		TimeStamp: currDate,
	}

	// Update bidHistory with new bid
	bidHistory = append(bidHistory, newBid)

	// Update Bids(foundAuctionItem) with new bid
	foundAuctionItem.Bids = append(foundAuctionItem.Bids, newBid)

	log.Println("bid history:", bidHistory)

	c.JSON(http.StatusCreated, newBid)
}

// --------------------- WebSocket Handler ---------------------

func websocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Register client
	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	// Send initial data
	initMsg := map[string]interface{}{
		"type":  "INITIAL_DATA",
		"items": auctionItems,
	}
	conn.WriteJSON(initMsg)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		var wsMsg WSMessage
		if err := json.Unmarshal(msg, &wsMsg); err != nil {
			conn.WriteJSON(gin.H{"type": "ERROR", "message": "Invalid message format"})
			continue
		}

		if wsMsg.Type == "NEW_BID" {
			handleNewBid(conn, wsMsg)
		} else {
			conn.WriteJSON(gin.H{"type": "ERROR", "message": "Invalid message type"})
			continue
		}
	}

	// Unregister on disconnect
	clientsMutex.Lock()
	delete(clients, conn)
	clientsMutex.Unlock()
	fmt.Println("Client disconnected")
}

// --------------------- Handle NEW_BID ---------------------

func handleNewBid(conn *websocket.Conn, data WSMessage) {
	if data.ItemId == 0 || data.BidAmount == 0 || len(data.Bidder) == 0 {
		conn.WriteJSON(gin.H{"type": "ERROR", "message": "Missing required fields"})
		return
	}

	foundAuctionItem := searchForAuctionItemById(data.ItemId)

	if foundAuctionItem == nil {
		conn.WriteJSON(gin.H{"type": "ERROR", "message": "Item not found"})
		return
	}

	if data.BidAmount <= foundAuctionItem.CurrentBid {
		conn.WriteJSON(gin.H{"type": "ERROR", "message": "Bid must be higher than current bid"})
		return
	}

	// Update CurrentBid with new bid amount
	foundAuctionItem.CurrentBid = data.BidAmount

	newId := len(bidHistory) + 1
	currDate := time.Now().UTC().Format(time.RFC3339)

	// Create new bid
	newBid := BidItem{
		Id:        uint32(newId),
		ItemId:    data.ItemId,
		Bidder:    data.Bidder,
		Amount:    data.BidAmount,
		TimeStamp: currDate,
	}

	// Update bidHistory with new bid
	bidHistory = append(bidHistory, newBid)

	// Update Bids(foundAuctionItem) with new bid
	foundAuctionItem.Bids = append(foundAuctionItem.Bids, newBid)

	log.Println("bid history:", bidHistory)

	// Broadcast update
	broadcastBidUpdate(foundAuctionItem, newBid)
}

// --------------------- Broadcast to All Clients ---------------------

func broadcastBidUpdate(item *AuctionItem, bid BidItem) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	msg := map[string]interface{}{
		"type": "BID_UPDATE",
		"item": item,
		"bid":  bid,
	}

	for client := range clients {
		if err := client.WriteJSON(msg); err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

// --------------------- WebSocket End ---------------------

func searchForAuctionItemById(id uint32) *AuctionItem {
	/* Lets find the auctionItem in the array using id
	   then reference it by use of foundAuctionItem
	*/

	var foundAuctionItem *AuctionItem

	for i := range auctionItems {
		if id == auctionItems[i].Id {
			foundAuctionItem = &auctionItems[i]
			break
		}
	}

	return foundAuctionItem

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
	//router.GET("/", index)
	router.GET("/api/items", fetchAuctionItems)
	router.GET("/api/items/:id", fetchAuctionItemById)
	router.GET("/api/history", fetchBidHistory)
	router.POST("/api/bids", processAuctionBids)
	// WebSocket endpoint
	//router.GET("/ws", websocketHandler)
	router.GET("/", func(c *gin.Context) {
		if websocket.IsWebSocketUpgrade(c.Request) {
			// Handle as WebSocket
			websocketHandler(c)
		} else {
			// Handle as HTTP
			index(c)
		}
	})

	router.Run(serverAddr)
}
