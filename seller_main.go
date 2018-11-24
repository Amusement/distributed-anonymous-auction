package main

import "./seller"

import (
	"log"
	"os"
)

func main() {
	log.Println("Starting seller client")

	if len(os.Args) != 3 {
		log.Fatalf("Usage: seller_main.go [REST_address] [initial_config_file_location]")
		os.Exit(1)
	}
	seller.Initialize(os.Args[1], os.Args[2])
}
