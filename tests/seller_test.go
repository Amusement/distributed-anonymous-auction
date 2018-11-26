package tests

import (
	"../seller"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestGetItem(t *testing.T) {
	seller := seller.Initialize("test_config.json")
	go seller.StartAuction("127.0.0.1:8080")
	resp, _ := http.Get("http://localhost:8080/seller/Item")
	var item string
	json.NewDecoder(resp.Body).Decode(&item)
	if item != "Fancy chocolate" {
		t.Errorf("Item was incorrect, got: %d, want: %d.", item, "Fancy chocolate")
	}
}

func TestGetAuctioneers(t *testing.T) {
	seller := seller.Initialize("test_config.json")
	go seller.StartAuction("127.0.0.1:8181")
	resp, _ := http.Get("http://localhost:8181/seller/Auctioneers")
	var auctioneers []string
	json.NewDecoder(resp.Body).Decode(&auctioneers)
	if auctioneers[0] != "127.0.0.1:8081" || auctioneers[1] != "127.0.0.1:8082" {
		t.Errorf("Auctioneers were incorrect, got: %d, want: %d.", auctioneers, []string{"127.0.0.1:8081", "127.0.0.1:8082"})
	}
}

func TestGetRoundNumber(t *testing.T) {
	seller := seller.Initialize("test_config.json")
	go seller.StartAuction("127.0.0.1:8282")
	resp, _ := http.Get("http://localhost:8282/seller/round")
	var roundNumber int
	json.NewDecoder(resp.Body).Decode(&roundNumber)
	if roundNumber != 1 {
		t.Errorf("Item was incorrect, got: %d, want: %d.", roundNumber, 1)
	}
}

func TestGetPrices(t *testing.T) {
	seller := seller.Initialize("test_config.json")
	go seller.StartAuction("127.0.0.1:8383")
	resp, _ := http.Get("http://localhost:8383/seller/prices")
	var prices []int
	json.NewDecoder(resp.Body).Decode(&prices)
	fmt.Println(prices)
	if prices[0] != 300 || prices[1] != 400 || prices[2] != 500 {
		t.Errorf("Item was incorrect, got: %d, want: %d.", prices, []int{300, 400, 500})
	}
}
//
//func TestGetPublicKey(t *testing.T) {
//	seller := seller.Initialize("test_config.json")
//	go seller.StartAuction("127.0.0.1:8080")
//	resp, _ := http.Get("http://localhost:8080/seller/Item")
//	var item string
//	json.NewDecoder(resp.Body).Decode(&item)
//	if item != "Fancy chocolate" {
//		t.Errorf("Item was incorrect, got: %d, want: %d.", item, "Fancy chocolate")
//	}
//}
