package main

import (
    "./bidder"
    "log"
    "os"
    "fmt"
	"bufio"
	"strconv"
	"strings"
	"net"
)

// Hack to get current IP address
func thisIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func main() {
    //log.Println("Bidder client starting.")
    if len(os.Args) != 2 {
        log.Fatalf("Usage: bidder_client.go [seller_ip_address]")
        os.Exit(1)
    }

    // Initialize bidder
    bidder := bidder.InitBidder(os.Args[1], thisIP().String())

    fmt.Printf("The seller is selling \"%v\" at the following prices: %v.\n", bidder.RoundInfo.Item, bidder.RoundInfo.Prices)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your maximum bid: ")
	text, _ := reader.ReadString('\n')

	maxBid, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		fmt.Println("Your bid was not understood: ", err)
		os.Exit(1)
	}
	bidder.ProcessBid(maxBid)
	bidder.ListenSeller()
}
