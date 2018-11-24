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
       - Communication to auctioneers using their REST API
*/

import (
	"../common"
	"crypto/rsa"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var seller Seller

type Config struct {
	Auctioneers []string
	StartTime   string
}

type Seller struct {
	config     Config
	router     *mux.Router
	prices     []int
	currRound  int
	publicKey  rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func Initialize(address, configFile string) {
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

	// Set callbacks for REST
	rtr := mux.NewRouter()
	rtr.HandleFunc("/seller/item", GetItem).Methods("GET")
	rtr.HandleFunc("/seller/prices", GetPrices).Methods("GET")
	rtr.HandleFunc("/seller/auctioneers", GetAuctioneers).Methods("GET")
	rtr.HandleFunc("/seller/round", GetRoundNumber).Methods("GET")
	rtr.HandleFunc("/seller/startTime", GetStartTime).Methods("GET")
	rtr.HandleFunc("/seller/timeLimit", GetTimeLimit).Methods("GET")
	rtr.HandleFunc("/seller/key", GetPublicKey).Methods("GET")

	// Create a global seller
	privK, pubK := common.GenerateRSA() // Generate RSA key pair
	seller = Seller{
		config:     config,
		router:     rtr,
		publicKey:  pubK,
		privateKey: privK,
	}

	// Run the REST server
	log.Println("Starting the seller server...")
	err = http.ListenAndServe(address, rtr)
	log.Printf("Error: %v", err)
	log.Printf("Public key: %v", seller.publicKey)
}

func GetItem(w http.ResponseWriter, r *http.Request) {

}

func GetPrices(w http.ResponseWriter, r *http.Request) {

}

func GetAuctioneers(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(seller.config.Auctioneers)
	if err != nil {
		log.Fatalf("error on GetAuctioneers: %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

func GetRoundNumber(w http.ResponseWriter, r *http.Request) {

}

func GetStartTime(w http.ResponseWriter, r *http.Request) {

}

func GetTimeLimit(w http.ResponseWriter, r *http.Request) {

}

func GetPublicKey(w http.ResponseWriter, r *http.Request) {
	data := common.MarshalKeyToPem(seller.publicKey)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

// Seller's private function ===========

func decodeID(msg []byte) {
	// Attempt to decode the message. If the decoded message is not in ip + price, we go to next round
	msg, err := common.DecryptID(msg, seller.privateKey)
	if err != nil {
		log.Fatalf("Error decrypting message: %v", err)
		// handle error
	}
	log.Printf("decoded msg: %v", string(msg))
}
