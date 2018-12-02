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
	"time"
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
	SellerIpPort   string
	LocalIpPort    string
	ExternalIpPort string
}

type AuctionRpcServer struct {
	Auctioneer *Auctioneer
}

func Initialize(config Config) *Auctioneer {
	return &Auctioneer{config: config,
		currentBids: make(map[common.Price][]common.Point),
		bidders:     make(map[string]struct{}),
		bidMutex:    &sync.Mutex{}}
}

func (a *Auctioneer) Start() {
	a.UpdateRoundInfo()

	rtr := mux.NewRouter()
	rtr.HandleFunc("/auctioneer/sendBid", a.SendBid).Methods("POST")
	rtr.HandleFunc("/auctioneer/compressedPoints", a.GetCompressedPoints).Methods("GET")
	go a.runAuction()
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
	status := a.roundInfo.AuctionStatus()
	if status == common.DURING {

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
	} else {
		w.WriteHeader(503)
	}
}

func (a *Auctioneer) GetCompressedPoints(w http.ResponseWriter, r *http.Request) {
	a.bidMutex.Lock()
	defer a.bidMutex.Unlock()
	err := json.NewEncoder(w).Encode(a.calculateCompressedPoints())
	if err != nil {
		fmt.Println(err)
	}
}

func (a *Auctioneer) calculateCompressedPoints() common.CompressedPoints {
	compressedPoints := common.CompressedPoints{make( map[common.Price]common.Point)}
	for key, points := range a.currentBids {
		var sum big.Int
		for _, p := range points {
			point := p.Y.Val
			sum.Add(&sum, point)
		}
		compressedPoints.Points[key] = common.Point{points[0].X, common.BigInt{&sum}}
	}

	return compressedPoints
}

//Meant to run a go routine
func (a *Auctioneer) runAuction() {
	until := time.Until(a.roundInfo.StartTime)
	fmt.Println("Waiting for ", until, " before auction")
	time.Sleep(until)
	fmt.Println("Auction started")
	time.Sleep(a.roundInfo.Interval.Duration)
	fmt.Println("Collecting compressed points")

	var compressedPoints []common.CompressedPoints
	//TODO: Add our own compressed points?
	for _, ipPort := range a.roundInfo.Auctioneers {
		if ipPort != a.config.ExternalIpPort{
			points, err := a.QueryCompressed(ipPort)
			if err == nil {
				compressedPoints = append(compressedPoints, points)
			} else {
				fmt.Println(err)
			}
		}
	}
	fmt.Println("Compressed points received", compressedPoints)
	res := common.ComputeLagrange(compressedPoints)
	fmt.Println(res)
	// res is map[Price]*big.Int, a map containing lagrange interpolation of each respective price
	// Making an assumption that T value is equal to number of auctioneers, will soon accept a T value
}

func (a *Auctioneer) QueryCompressed(ipPort string) (common.CompressedPoints, error) {
	fmt.Println("Getting compressed point from ", ipPort)
	req, err := http.NewRequest("GET", "http://"+ipPort+"/auctioneer/compressedPoints", nil)
	client := &http.Client{}
	var compressedPoints common.CompressedPoints

	resp, err := client.Do(req)
	if err != nil {
		return compressedPoints, err
	}

	defer resp.Body.Close()

	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	if err := json.NewDecoder(resp.Body).Decode(&compressedPoints); err != nil {
		return compressedPoints, err
	}

	return compressedPoints, nil
}
