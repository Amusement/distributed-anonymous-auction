/*
This package specifies the seller API
*/

package seller

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var seller Seller

type Seller struct {
	router   *mux.Router
	prices   []int
}

func Initialize(initialPrices []int) (err error) {
	fmt.Println("Seller initialized")
	rtr := mux.NewRouter()
	seller = Seller{router: rtr, prices: initialPrices}
	rtr.HandleFunc("/seller/item", GetItem).Methods("GET")
	rtr.HandleFunc("/seller/prices", GetPrices).Methods("GET")
	rtr.HandleFunc("/seller/auctioneers", GetAuctioneers).Methods("GET")
	rtr.HandleFunc("/seller/round", GetRoundNumber).Methods("GET")
	rtr.HandleFunc("/seller/startTime", GetStartTime).Methods("GET")
	rtr.HandleFunc("/seller/timeLimit", GetTimeLimit).Methods("GET")
	rtr.HandleFunc("/seller/key", GetPublicKey).Methods("GET")
	return http.ListenAndServe(":8000", rtr)
}

func GetItem(w http.ResponseWriter, r *http.Request) {

}

func GetPrices(w http.ResponseWriter, r *http.Request) {

}

func GetAuctioneers(w http.ResponseWriter, r *http.Request) {

}

func GetRoundNumber(w http.ResponseWriter, r *http.Request) {

}

func GetStartTime(w http.ResponseWriter, r *http.Request) {

}

func GetTimeLimit(w http.ResponseWriter, r *http.Request) {

}

func GetPublicKey(w http.ResponseWriter, r *http.Request) {

}
