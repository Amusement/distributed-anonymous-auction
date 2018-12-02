package seller

/*
   Currently working feature
       - Seller can generate its own RSA key value
       - Basic REST api function working
           - Can query public key
           - Can query auctioneer list
           - Can query  prices
           - Can query current round
           - Can query start time

           - Can query t_value

   Need to implement (not a whole list)
       - Prices logic
       - Round logic
           - Start time / End time logic
       - Winner declaration logic
           - decodeID function is implemented, use this to figure it out
       - Communication to Auctioneers using their REST API
*/

import (
	"../common"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// TODO: Consider moving some fields into AuctionRound type

type Seller struct {
	AuctionRound          common.AuctionRound
	waitingForCalculation bool
	auctionIsOver         bool
	router                *mux.Router
	publicKey             rsa.PublicKey
	privateKey            *rsa.PrivateKey
}

func Initialize(configFile string) *Seller {
	// Get configuration of the seller
	var auctionRound common.AuctionRound
	file, err := os.Open(configFile)
	defer file.Close()
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
		os.Exit(1)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&auctionRound)
	if err != nil {
		log.Fatalf("Error decoding config file: %v", err)
		os.Exit(1)
	}

	// Create a new router
	rtr := mux.NewRouter()

	// Create a global seller
	privK, pubK := common.GenerateRSA() // Generate RSA key pair
	seller := &Seller{
		AuctionRound:          auctionRound,
		waitingForCalculation: false,
		auctionIsOver:         false,
		router:                rtr,
		publicKey:             pubK,
		privateKey:            privK,
	}
	return seller
}

func (s *Seller) checkRoundTermination() {
	timeForEnd := time.Until(s.AuctionRound.StartTime.Add(s.AuctionRound.Interval.Duration))
	time.Sleep(timeForEnd)
	s.waitingForCalculation = true
	// TODO: Receive ids of highest price range from auctioneers
	// TODO: Set waiting for calculation to false and determine if there will be another round or auction is over
}

func (s *Seller) StartAuction(address string) {
	s.router.HandleFunc("/seller/key", s.GetPublicKey).Methods("GET")
	s.router.HandleFunc("/seller/roundinfo", s.GetRoundInfo).Methods("GET")
	s.router.HandleFunc("/seller/auctionover", s.GetAuctionOverStatus).Methods("GET")
	s.router.HandleFunc("/seller/waitingcalculation", s.GetWaitingCalculationStatus).Methods("GET")
	// Run the REST server
	go s.checkRoundTermination()
	log.Printf("Error: %v", http.ListenAndServe(address, s.router))
	// TODO remove this sleep after
	time.Sleep(10000 * time.Second)
}

func (s *Seller) GetPublicKey(w http.ResponseWriter, r *http.Request) {
	data := common.MarshalKeyToPem(s.publicKey)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func (s *Seller) GetRoundInfo(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(s.AuctionRound)
	if err != nil {
		log.Fatalf("error on GetRoundInfo: %v", err)
	}
	w.Write(data)
}

func (s *Seller) GetAuctionOverStatus(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(s.auctionIsOver)
	if err != nil {
		log.Fatalf("error on GetAuctioOverStatus: %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func (s *Seller) GetWaitingCalculationStatus(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(s.waitingForCalculation)
	if err != nil {
		log.Fatalf("error on GetWaitingCalculationStatus: %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

// Seller's private function ===========

func (s *Seller) decodeID(msg []byte) {
	// Attempt to decode the message. If the decoded message is not in ip + price, we go to next round
	msg, err := common.DecryptID(msg, s.privateKey)
	if err != nil {
		log.Fatalf("Error decrypting message: %v", err)
		// handle error
	}
	log.Printf("decoded msg: %v", string(msg))
}

func (s *Seller) contactWinner(ipPortAndPrice string) {
	ipPort := strings.Split(ipPortAndPrice, " ")[0]
	conn, err := net.Dial("tcp", ipPort)
	defer conn.Close()
	if err != nil {
		fmt.Println("Was not able to contact winning bidder: ", err)
	}
	conn.Write([]byte("winner"))
}

func (s *Seller) calculateNewRound(highestBid uint) {
	prices, _ := s.CalculateNewPrices(highestBid)
	newAuctionRound := common.AuctionRound{
		Item:         s.AuctionRound.Item,
		StartTime:    time.Now(),
		Interval:     s.AuctionRound.Interval,
		Prices:       prices,
		Auctioneers:  s.AuctionRound.Auctioneers,
		T:            s.AuctionRound.T,
		CurrentRound: s.AuctionRound.CurrentRound + 1,
	}
	s.AuctionRound = newAuctionRound

}

func (s *Seller) CalculateNewPrices(highestBid uint) ([]uint, error) {
	numberOfPrices := len(s.AuctionRound.Prices)
	if numberOfPrices != 0 {
		priceInteval := s.AuctionRound.Prices[1] - s.AuctionRound.Prices[0]
		if s.AuctionRound.Prices[numberOfPrices-1] == highestBid {
			var newPrices []uint
			for i := 0; i < numberOfPrices; i++ {
				newPrices = append(newPrices, highestBid+uint(i)*priceInteval)
			}
			return newPrices, nil
		} else {
			newPriceInterval := uint(math.Ceil(float64(priceInteval) / float64(numberOfPrices)))
			var newPrices []uint
			for i := 0; i < numberOfPrices; i++ {
				newPrices = append(newPrices, uint(highestBid + uint(i)*newPriceInterval))
			}
			return newPrices, nil
		}
	} else {
		return nil, errors.New("Seller price list is empty!")
	}
}
