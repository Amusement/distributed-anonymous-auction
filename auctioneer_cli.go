package main

import "./auctioneer"

func main() {
	auctioneer := auctioneer.Initialize(auctioneer.Config{"127.0.0.1:6000", "127.0.0.1"})
	auctioneer.Start()
}
