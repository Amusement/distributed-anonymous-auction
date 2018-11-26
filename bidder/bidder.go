package bidder

/*
   Currently working feature
       - Basic initializtion of bidder
           - Can get list of auctioneers, public key, startTime, prices from seller

   Need to implement (not a whole list)
       - Round logic helper fuctions for bidder_client.go

*/

import (
	"encoding/json"
	"github.com/jongukim/polynomial"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
)

// TODO: Finalize types
type Bidder struct {
	// List of clients
	listOfAuctioneers []string
	publicKey         string
	sellerIP          string
	startTime         string
	prices		  []string
}

func InitBidder(sellerAddr string) Bidder {
	//var rawConfig string

	// unmarshal rawConfig, return a bidder
	b := &Bidder{
		sellerIP: sellerAddr,
	}
	b.getConfig()
	log.Printf("bidder: %v", b)
	return *b
}

// =============== REST call to seller ===============

// Get configuration from the seller
func (b *Bidder) getConfig() {
	// Get public key
	uri := b.sellerIP + "/seller/key"
	response, err := http.Get(uri)
	if err != nil {
		log.Fatalf("Failed to get config file from seller: %v", err)
		os.Exit(1)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		b.publicKey = string(data)
	}

	// Get List of auctioneers
	uri = b.sellerIP + "/seller/auctioneers"
	response, err = http.Get(uri)
	if err != nil {
		log.Fatalf("Failed to get config file from seller: %v", err)
		os.Exit(1)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(data, &b.listOfAuctioneers)
	}

	// Get start time
	uri = b.sellerIP + "/seller/time/start"
	response, err = http.Get(uri)
	if err != nil {
		log.Fatalf("Failed to get starting time from seller: %v", err)
		os.Exit(1)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		b.startTime = string(data)
	}

	// Get prices 
	uri = b.sellerIP + "/seller/prices"
	response, err = http.Get(uri)
	if err != nil {
		log.Fatalf("Failed to get prices from seller: %v", err)
		os.Exit(1)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(data, &b.prices)
	}

	// Generate N polynomials depending on price range
}

// ============== Other functions ==================
func generatePolynomial(degree int, id *big.Int) []*big.Int {
	// f(x) = 3x^3 + 2x + 1 => [1 2 0 3]

	poly := polynomial.RandomPoly(int64(degree), 5) // 5 is hard coded to make coefficients 2^5 at most

	// Change the ID
	poly[0] = id
	return poly
}
