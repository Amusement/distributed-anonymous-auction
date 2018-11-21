package auctioneer

import (
	"../common"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
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
	server := &AuctionRpcServer{a}
	rpc.Register(server)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", a.config.LocalIpPort)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

// Receives bids from a bidder and returns if true if it was successfully received
func (a *AuctionRpcServer) ReceiveBid(bidPoints *common.BidPoints, reply *bool) error {
	a.Auctioneer.bidMutex.Lock()
	a.Auctioneer.currentBids[bidPoints.BidderID] = bidPoints.Points
	a.Auctioneer.bidMutex.Unlock()
	*reply = true
	return nil

}
