package main

import (
	"./seller"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.Println("Starting seller client")
	if len(os.Args) != 3 {
		log.Fatalf("Usage: seller_main.go [REST_address (IP:PORT)] [initial_config_file_location]")
		os.Exit(1)
	}
	s := seller.Initialize(os.Args[2])
	go s.StartAuction(os.Args[1])
	log.Println("Started seller REST API")

	// Start the auction CLI
	if time.Now().UTC().After(s.AuctionRound.StartTime) {
		log.Printf("start time: %v", s.AuctionRound.StartTime)
		log.Printf("now: %v", time.Now().UTC())
		fmt.Println("Invalid start time.")
		return
	}

	fmt.Println("\n\n=====Starting the auction!=====")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Enter list of prices separated by comma for round %v:\n", s.AuctionRound.CurrentRound)
		priceRange, _ := reader.ReadString('\n')
		//priceStrings := strings.Split(strings.TrimSpace(priceRange), "\n")
		priceStrings := strings.Split(strings.TrimSpace(priceRange), ",")
		var prices = []uint{}

		log.Println("Got pricerange: ", priceStrings, " of length: ", len(priceRange))

		for _, i := range priceStrings {
			j, err := strconv.ParseUint(i, 10, 32)
			if err != nil {
				panic(err)
			}
			prices = append(prices, uint(j))
		}

		s.AuctionRound.Prices = prices
		s.AuctionRound.CurrentRound += 1

		fmt.Println("Waiting for current round to finish...")
		for s.AuctionRound.StartTime.Add(s.AuctionRound.Interval.Duration).After(time.Now().UTC()) {
			time.Sleep(time.Second) // Sleep and check every 1 second
		}

		// Ask auctioneers for points
		// Check if there is a winner
		//   -- if winner, contact him
		//   -- if tie, compute new round and start new round
		// Check for winner
		// Traverse through auctioneer list and call its REST API
		// TODO

	}
}
