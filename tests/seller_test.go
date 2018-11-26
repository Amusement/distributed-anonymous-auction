package tests

import (
	"../seller"
	"encoding/json"
	"net/http"
	"testing"
	"time"
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
	var prices []string
	json.NewDecoder(resp.Body).Decode(&prices)
	if prices[0] != "300" || prices[1] != "400" || prices[2] != "500" {
		t.Errorf("Item was incorrect, got: %d, want: %d.", prices, []string{"300", "400", "500"})
	}
}

func TestGetTvalue(t *testing.T) {
	seller := seller.Initialize("test_config.json")
	go seller.StartAuction("127.0.0.1:8484")
	resp, _ := http.Get("http://localhost:8484/seller/tvalue")
	var tvalue int
	json.NewDecoder(resp.Body).Decode(&tvalue)
	if tvalue != 2 {
		t.Errorf("Item was incorrect, got: %d, want: %d.", tvalue, 2)
	}
}

func TestGetStartTime(t *testing.T) {
	seller := seller.Initialize("test_config.json")
	go seller.StartAuction("127.0.0.1:8585")
	resp, _ := http.Get("http://localhost:8585/seller/time/start")
	var startTime time.Time
	json.NewDecoder(resp.Body).Decode(&startTime)
	expectedTime, _ := time.Parse("2006-01-02T15:04:05Z07:00","2018-11-28T16:30:00Z")
	if startTime.String() != expectedTime.String() {
		t.Errorf("Item was incorrect, got: %d, want: %d.", startTime.String(), expectedTime.String())
	}
}

func TestGetTimeLimit(t *testing.T) {
	seller := seller.Initialize("test_config.json")
	go seller.StartAuction("127.0.0.1:8686")
	resp, _ := http.Get("http://localhost:8686/seller/time/limit")
	var limit int
	json.NewDecoder(resp.Body).Decode(&limit)
	if limit != 60 {
		t.Errorf("Item was incorrect, got: %d, want: %d.", limit, 60)
	}
}

func TestGetTimeInterval(t *testing.T) {
	seller := seller.Initialize("test_config.json")
	go seller.StartAuction("127.0.0.1:8787")
	resp, _ := http.Get("http://localhost:8787/seller/time/interval")
	var interval int
	json.NewDecoder(resp.Body).Decode(&interval)
	if interval != 60000000000 {
		t.Errorf("Item was incorrect, got: %d, want: %d.", interval, 60000000000)
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
