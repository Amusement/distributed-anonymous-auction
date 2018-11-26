package main

import (
    "./bidder"
    "log"
    "os"
)

func main() {
    log.Println("Bidder client starting.")

    if len(os.Args) != 2 {
        log.Fatalf("Usage: bidder_client.go [seller_ip_address]")
        os.Exit(1)
    }

    // Initialize a bidder
    //bidder := bidder.InitBidder(os.Args[1])
    bidder.InitBidder(os.Args[1])

}
