package auctioneer

import (
	"../common"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/big"
	"net/http"
	"sync"
)

type Auctioneer struct {
	config      Config
	round       uint
	bidMutex    *sync.Mutex
	currentBids map[common.Price][]common.Point
	bidders     map[string]struct{}
	roundInfo   common.AuctionRound
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
		currentBids: make(map[common.Price][]common.Point),
		bidders:     make(map[string]struct{}),
		bidMutex:    &sync.Mutex{}}
}

func (a *Auctioneer) Start() {
	a.UpdateRoundInfo()

	rtr := mux.NewRouter()
	rtr.HandleFunc("/auctioneer/sendBid", a.SendBid).Methods("POST")
	rtr.HandleFunc("/auctioneer/compressedPoints", a.GetCompressedPoints).Methods("GET")

	log.Println("Starting the auctioneer server...")
	err := http.ListenAndServe(a.config.LocalIpPort, rtr)
	log.Printf("Error: %v", err)
}

func (a *Auctioneer) UpdateRoundInfo() {
	req, err := http.NewRequest("GET", "http://"+a.config.SellerIpPort+"/seller/roundinfo", nil)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	defer resp.Body.Close()

	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	var roundInfo common.AuctionRound
	if err := json.NewDecoder(resp.Body).Decode(&roundInfo); err != nil {
		log.Println(err)
	}
	a.roundInfo = roundInfo
}

// Receives bids from a bidder and returns if true if it was successfully received

func (a *Auctioneer) SendBid(w http.ResponseWriter, r *http.Request) {
	var bidPoints common.BidPoints
	err := json.NewDecoder(r.Body).Decode(&bidPoints)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Received bid from ", bidPoints.BidderID)
	a.bidMutex.Lock()
	for price, points := range bidPoints.Points {
		a.currentBids[price] = append(a.currentBids[price], points)
	}
	a.bidders[bidPoints.BidderID] = struct{}{}
	a.bidMutex.Unlock()
	w.WriteHeader(200)
}

func (a *Auctioneer) GetCompressedPoints(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(a.calculateCompressedPoints())
	if err != nil {
		fmt.Println(err)
	}
}

func (a *Auctioneer) calculateCompressedPoints() common.CompressedPoints {
	a.bidMutex.Lock()
	defer a.bidMutex.Unlock()
	compressedPoints := common.CompressedPoints{make(map[common.Price]common.Point)}
	for key, points := range a.currentBids {
		var sum big.Int
		for _, p := range points {
			point := p.Y.Val
			sum.Add(&sum, point)
		}
		compressedPoints.Points[key] = common.Point{points[0].X, common.BigInt{&sum}}
	}
	fmt.Println("Price point ", compressedPoints)

	return compressedPoints
}
