package main

import "./auctioneer"

func main() {
	auctioneer := auctioneer.Initialize(auctioneer.Config{"127.0.0.1:5555", "127.0.0.1"})
	auctioneer.Start()
}
