package seller

/*
   Currently working feature
       - Seller can generate its own RSA key value
       - Basic REST api function working
           - Can query public key
           - Can query auctioneer list

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

type Config struct {
	Item        string
	Prices      []int
	RoundNumber int
	Auctioneers []string
	StartTime   time.Time
	TimeLimit   int
	Interval    time.Duration
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
	s.router.HandleFunc("/seller/key", s.GetPublicKey).Methods("GET")
	s.router.HandleFunc("/seller/Auctioneers", s.GetAuctioneers).Methods("GET")
	s.router.HandleFunc("/seller/round", s.GetRoundNumber).Methods("GET")
	s.router.HandleFunc("/seller/prices", s.GetPrices).Methods("GET")
	s.router.HandleFunc("/seller/Item", s.GetItem).Methods("GET")
	// TODO: Add more functions
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
	data, _ := json.Marshal(s.Config.RoundNumber)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func (s *Seller) GetPrices(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(s.Config.Prices)
	if err != nil {
		log.Fatalf("error on GetAuctioneers: %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func (s *Seller) GetItem(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(s.Config.Item)
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
