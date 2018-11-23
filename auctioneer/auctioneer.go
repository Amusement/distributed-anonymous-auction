package auctioneer

import (
	"../common"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"

	"github.com/googollee/go-socket.io"
)

type Auctioneer struct {
	config      Config
	round       uint
	bidMutex    *sync.Mutex
	currentBids map[string][]common.Point
}

type Config struct {
	LocalIpPort string
	ExternalIp  string
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

func (a *Auctioneer) Start() {
	socketSever, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	socketSever.On("connection", func(so socketio.Socket) {
		so.On("bidder", func(msg string) {
			log.Println(so.Id(), " got bid from ", msg)
		})
		so.On("disconnection", func() {
			log.Println("Disconnected from peer")
		})
	})

	socketSever.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", socketSever)

	rtr := mux.NewRouter()
	rtr.HandleFunc("/auctioneer/sendBid", a.SendBid).Methods("POST")

	// Run the REST server
	log.Println("Starting the auctioneer server...")
	err = http.ListenAndServe(a.config.LocalIpPort, rtr)
	log.Printf("Error: %v", err)
}

// Receives bids from a bidder and returns if true if it was successfully received

func (a *Auctioneer) SendBid(w http.ResponseWriter, r *http.Request) {
	var bidPoints common.BidPoints
	_ = json.NewDecoder(r.Body).Decode(&bidPoints)

	fmt.Println("Received bid from ", bidPoints.BidderID)
	a.bidMutex.Lock()
	a.currentBids[bidPoints.BidderID] = bidPoints.Points
	a.bidMutex.Unlock()
	w.WriteHeader(200)
}
