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
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

// TODO: Consider moving some fields into AuctionRound type
type Config struct {
	Item        string
	Prices      []uint
	CurrRound   int
	Auctioneers []string
	StartTime   time.Time
	TimeLimit   int
	Interval    common.Duration
	T_value     int
}

type Seller struct {
	Config     Config
	router     *mux.Router
	publicKey  rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func Initialize(configFile string) *Seller {
	// Get configuration of the seller
	var config Config
	file, err := os.Open(configFile)
	defer file.Close()
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
		os.Exit(1)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Error decoding config file: %v", err)
		os.Exit(1)
	}

	// Create a new router
	rtr := mux.NewRouter()

	// Create a global seller
	privK, pubK := common.GenerateRSA() // Generate RSA key pair
	seller := &Seller{
		Config:     config,
		router:     rtr,
		publicKey:  pubK,
		privateKey: privK,
	}
	return seller
}

func (s *Seller) StartAuction(address string) {
	s.router.HandleFunc("/seller/roundinfo", s.GetRoundInfo).Methods("GET")
	s.router.HandleFunc("/seller/key", s.GetPublicKey).Methods("GET")
	s.router.HandleFunc("/seller/Auctioneers", s.GetAuctioneers).Methods("GET")
	s.router.HandleFunc("/seller/round", s.GetRoundNumber).Methods("GET")
	s.router.HandleFunc("/seller/prices", s.GetPrices).Methods("GET")
	s.router.HandleFunc("/seller/Item", s.GetItem).Methods("GET")
	s.router.HandleFunc("/seller/tvalue", s.GetTValue).Methods("GET")
	s.router.HandleFunc("/seller/time/start", s.GetStartTime).Methods("GET")
	s.router.HandleFunc("/seller/time/limit", s.GetTimeLimit).Methods("GET")
	s.router.HandleFunc("/seller/time/interval", s.GetInterval).Methods("GET")
	// Run the REST server
	log.Printf("Error: %v", http.ListenAndServe(address, s.router))
}

func (s *Seller) GetPublicKey(w http.ResponseWriter, r *http.Request) {
	data := common.MarshalKeyToPem(s.publicKey)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func (s *Seller) GetAuctioneers(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(s.Config.Auctioneers)
	if err != nil {
		log.Fatalf("error on GetAuctioneers: %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}


func (s *Seller) GetRoundNumber(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(s.Config.CurrRound)
	if err != nil {
		log.Fatalf("error on GetRoundNumber: %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func (s *Seller) GetPrices(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(s.Config.Prices)
	if err != nil {
		log.Fatalf("error on GetPrices: %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func (s *Seller) GetItem(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(s.Config.Item)
	if err != nil {
		log.Fatalf("error on GetItem: %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func (s *Seller) GetTValue(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(s.Config.T_value)
	if err != nil {
		log.Fatalf("error on GetTvalue: %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func (s *Seller) GetStartTime(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(s.Config.StartTime)
	if err != nil {
		log.Fatalf("error on GetStartTime: %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func (s *Seller) GetTimeLimit(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(s.Config.TimeLimit)
	if err != nil {
		log.Fatalf("error on GetTimeLimit: %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func (s *Seller) GetInterval(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(s.Config.Interval)
	if err != nil {
		log.Fatalf("error on GetInterval: %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}


func (s *Seller) GetRoundInfo(w http.ResponseWriter, r *http.Request) {
	convertedRoundInfo := common.AuctionRound{
		s.Config.Item,
		s.Config.StartTime,
		s.Config.Interval,
		s.Config.Prices,
		s.Config.Auctioneers,
		s.Config.T_value,
		s.Config.CurrRound,
	}

	data, err := json.Marshal(convertedRoundInfo)
	if err != nil {
		log.Fatalf("error on GetRoundInfo: %v", err)
	}

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
