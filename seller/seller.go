/*
This package specifies the seller API
*/

package seller

import (
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
	PublicKey   string
}

type Seller struct {
	config    Config
	router    *mux.Router
	prices    []int
	currRound int
}

func Initialize(port, configFile string) {
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
	seller = Seller{
		config: config,
		router: rtr,
	}

	// Run the REST server
	log.Println("Starting the seller server...")
	err = http.ListenAndServe(":"+port, rtr)
	log.Printf("Error: %v", err)
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
	data, err := json.Marshal(seller.config.PublicKey)
	if err != nil {
		log.Fatalf("error on GetAuctioneers: %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}
