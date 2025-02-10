package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"github.com/google/uuid"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price           string `json:"price"`
}

type ResponseID struct {
	ID string `json:"id"`
}

type ResponsePoints struct {
	Points int `json:"points"`
}

var (
	receipts = make(map[string]int)
	mutex    = &sync.Mutex{}
)

func calculatePoints(receipt Receipt) int {
	points := 0

	// One point for every alphanumeric character in the retailer name
	alphaNum := regexp.MustCompile(`[^a-zA-Z0-9]`)
	cleanRetailer := alphaNum.ReplaceAllString(receipt.Retailer, "")
	points += len(cleanRetailer)

	// Convert total to float
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil {
		if total == math.Floor(total) {
			points += 50 // 50 points if total is a round number
		}
		if math.Mod(total, 0.25) == 0 {
			points += 25 // 25 points if total is multiple of 0.25
		}
	}

	// 5 points for every two items on the receipt
	points += (len(receipt.Items) / 2) * 5

	// Points based on item description length
	for _, item := range receipt.Items {
		desc := strings.TrimSpace(item.ShortDescription)
		if len(desc)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				points += int(math.Ceil(price * 0.2))
			}
		}
	}

	// 6 points if purchase date is odd
	dateParts := strings.Split(receipt.PurchaseDate, "-")
	if len(dateParts) == 3 {
		day, err := strconv.Atoi(dateParts[2])
		if err == nil && day%2 == 1 {
			points += 6
		}
	}

	// 10 points if purchase time is between 2:00pm and 4:00pm
	timeParts := strings.Split(receipt.PurchaseTime, ":")
	if len(timeParts) == 2 {
		hour, err1 := strconv.Atoi(timeParts[0])
		minute, err2 := strconv.Atoi(timeParts[1])
		if err1 == nil && err2 == nil {
			if hour == 14 || (hour == 15 && minute < 60) {
				points += 10
			}
		}
	}

	return points
}

func processReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	points := calculatePoints(receipt)

	mutex.Lock()
	receipts[id] = points
	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponseID{ID: id})
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/receipts/")
	id = strings.TrimSuffix(id, "/points")

	mutex.Lock()
	points, exists := receipts[id]
	mutex.Unlock()

	if !exists {
		http.Error(w, "No receipt found for that ID", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponsePoints{Points: points})
}

func main() {
	http.HandleFunc("/receipts/process", processReceipt)
	http.HandleFunc("/receipts/", getPoints)
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

