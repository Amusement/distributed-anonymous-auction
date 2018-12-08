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

const NOBID = "No Bid"
const MULTIPLEWINNERS = "Multiple Winners"

// TODO: Consider moving some fields into AuctionRound type

type Seller struct {
	AuctionRound          common.AuctionRound
	router                *mux.Router
	publicKey             rsa.PublicKey
	privateKey            *rsa.PrivateKey
	// Key is Ip Port of auctioneer
	BidPoints map[string]map[common.Price]common.BigInt
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

	// Some simple validations
	if auctionRound.T >= len(auctionRound.Auctioneers) {
		log.Fatalf("config file error: T value should be lower than length of auctioneers")
		os.Exit(1)
	}

	auctionRound.CurrentRound = 1

	// Create a new router
	rtr := mux.NewRouter()

	// Create a global seller
	privK, pubK := common.GenerateRSA() // Generate RSA key pair
	seller := &Seller{
		AuctionRound:          auctionRound,
		router:                rtr,
		publicKey:             pubK,
		privateKey:            privK,
		BidPoints:             make(map[string]map[common.Price]common.BigInt),
	}
	return seller
}

func (s *Seller) checkRoundTermination() {
	for {
		// Waiting for bidding round to end
		timeForEnd := time.Until(s.AuctionRound.StartTime.Add(s.AuctionRound.Interval.Duration))
		time.Sleep(timeForEnd)
		fmt.Println("Round is over. Waiting for lagrange calculation and all that stuff.")
		// Waiting for calculating round to end
		time.Sleep(s.AuctionRound.Interval.Duration/common.IntervalMultiple)

		// Calculate for a winner
		if len(s.BidPoints) < s.AuctionRound.T {
			fmt.Println("We lost more than T auctioneers :(")
		}
		for _, priceMap := range s.BidPoints {
			isDone := false
			for i := len(priceMap) - 1; i >= 0; i-- {
				price := common.Price(s.AuctionRound.Prices[i])
				encryptedID := priceMap[price]
				res := s.decodeID(encryptedID.Val.Bytes())
				if res == NOBID {
					fmt.Println("There are no bids for price: ", price)
				} else if res == MULTIPLEWINNERS {
					fmt.Println("There are multiple winners for price: ", price)
					s.calculateNewRound(uint(price))
					isDone = true
					s.BidPoints = make(map[string]map[common.Price]common.BigInt)
					break
				} else {
					fmt.Println("Got a winner: ", res)
					s.contactWinner(res, price)
					s.AuctionRound.CurrentRound = -1
					time.Sleep(5 * time.Second)
					return
				}
			}
			if isDone == true {
				break
			}
		}
	}
}

func (s *Seller) StartAuction(address string) {
	s.router.HandleFunc("/seller/key", s.GetPublicKey).Methods("GET")
	s.router.HandleFunc("/seller/roundinfo", s.GetRoundInfo).Methods("GET")
	s.router.HandleFunc("/seller/bidpoint", s.PostBidsPoint).Methods("POST")

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

// receives points from the auctioneers
func (s *Seller) PostBidsPoint(w http.ResponseWriter, r *http.Request) {
	var totalBids common.TotalBids
	err := json.NewDecoder(r.Body).Decode(&totalBids)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.BidPoints[totalBids.AuctioneerId] = common.ComputeLagrange(totalBids.Points)
	w.WriteHeader(200)
}

// Seller's private function ===========

func (s *Seller) decodeID(msg []byte) string {
	if len(msg) == 0 {
		return "No Bid"
	}
	// Attempt to decode the message.
	rawMsg, err := common.DecryptID(msg, s.privateKey)
	if err != nil {
		return "Multiple Winners"
	}
	return string(rawMsg)
}

func (s *Seller) contactWinner(ipPortAndPrice string, price common.Price) {
	ipPort := strings.Split(ipPortAndPrice, " ")[0]
	conn, err := net.Dial("tcp", ipPort)
	defer conn.Close()
	if err != nil {
		fmt.Println("Was not able to contact winning bidder: ", err)
	}
	winnerNotification := common.WinnerNotification{ WinningPrice: price }
	//encoder := gob.NewEncoder(conn)
	//encoder.Encode(winnerNotification)
	notifBytes, err := json.Marshal(winnerNotification)
	if err != nil {
		fmt.Println("Issue encoding winner notification: ", err)
	}
	conn.Write(notifBytes)
}

func (s *Seller) calculateNewRound(highestBid uint) {
	prices, _ := s.CalculateNewPrices(highestBid)
	newAuctionRound := common.AuctionRound{
		Item:         s.AuctionRound.Item,
		StartTime:    time.Now().Add(s.AuctionRound.Interval.Duration),
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
		priceInterval := s.AuctionRound.Prices[1] - s.AuctionRound.Prices[0]
		if s.AuctionRound.Prices[numberOfPrices-1] == highestBid {
			var newPrices []uint
			for i := 0; i < numberOfPrices; i++ {
				newPrices = append(newPrices, highestBid+uint(i)*priceInterval)
			}
			return newPrices, nil
		} else {
			newPriceInterval := uint(math.Ceil(float64(priceInterval) / float64(numberOfPrices)))
			var newPrices []uint
			for i := 0; i < numberOfPrices; i++ {
				newPrices = append(newPrices, uint(highestBid+uint(i)*newPriceInterval))
			}
			return newPrices, nil
		}
	} else {
		return nil, errors.New("Seller price list is empty!")
	}
}
