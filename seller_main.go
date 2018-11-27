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
	if time.Now().UTC().After(s.Config.StartTime) {
		log.Printf("start time: %v", s.Config.StartTime)
		log.Printf("now: %v", time.Now().UTC())
		fmt.Println("Invalid start time.")
		return
	}

	fmt.Println("\n\n=====Starting the auction!=====")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Enter list of prices separated by comma for round %v:\n", s.Config.CurrRound)
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

		s.Config.Prices = prices
		s.Config.CurrRound += 1

		fmt.Println("Waiting for current round to finish...")
		for s.Config.StartTime.Add(s.Config.Interval).After(time.Now().UTC()) {
			time.Sleep(time.Second) // Sleep and check every 1 second
		}

		// Check for winner
		// Traverse through auctioneer list and call its REST API
		// TODO

	}
}
