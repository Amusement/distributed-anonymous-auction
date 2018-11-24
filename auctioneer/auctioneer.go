package auctioneer

import (
	"../common"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rsms/gotalk"
	"log"
	"net/http"
	"sync"
)

type Auctioneer struct {
	config      Config
	round       uint
	bidMutex    *sync.Mutex
	currentBids map[string][]common.Point
	peers       []string
}

type Config struct {
	SellerIpPort string
	LocalIpPort  string
	ExternalIp   string
}

type AuctionRpcServer struct {
	Auctioneer *Auctioneer
}

func Initialize(config Config) *Auctioneer {
	return &Auctioneer{config: config,
		round:       0,
		currentBids: make(map[string][]common.Point),
		bidMutex:    &sync.Mutex{}}
}

type bidderMsg struct {
	source   string
	bidderID string
}

func (a *Auctioneer) Start() {
	gotalk.Handle("bidder", func(msg bidderMsg) (error) {
		fmt.Println(msg.source, " got a bid from ", msg.bidderID)
		return nil
	})
	http.Handle("/gotalk/", gotalk.WebSocketHandler())

	a.UpdatePeers()

	rtr := mux.NewRouter()
	rtr.HandleFunc("/auctioneer/sendBid", a.SendBid).Methods("POST")

	log.Println("Starting the auctioneer server...")
	err := http.ListenAndServe(a.config.LocalIpPort, rtr)
	log.Printf("Error: %v", err)
}

func (a *Auctioneer) UpdatePeers() {
	req, err := http.NewRequest("GET", "http://"+a.config.SellerIpPort+"/seller/auctioneers", nil)
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	var peers []string
	if err := json.NewDecoder(resp.Body).Decode(&peers); err != nil {
		log.Println(err)
	}
	a.peers = peers
}

// Receives bids from a bidder and returns if true if it was successfully received

func (a *Auctioneer) SendBid(w http.ResponseWriter, r *http.Request) {
	var bidPoints common.BidPoints
	_ = json.NewDecoder(r.Body).Decode(&bidPoints)

	gotalk.Connect("tcp", )
	fmt.Println("Received bid from ", bidPoints.BidderID)
	a.bidMutex.Lock()
	a.currentBids[bidPoints.BidderID] = bidPoints.Points
	a.bidMutex.Unlock()
	w.WriteHeader(200)
}
