package bidder

/*
   Currently working:
       - Basic initializtion of bidder
           - Can get list of auctioneers, public key, startTime, Prices from seller

*/

import (
	"encoding/json"
	"github.com/jongukim/polynomial"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"P2-d3w9a-b3c0b-b3l0b-k0b9/common"
	"fmt"
)

type Bidder struct {
	RoundInfo		common.AuctionRound			// Initially retrieved round info.
	sellerPublicKey string
	sellerIP        string
}

func InitBidder(sellerAddr string) *Bidder {
	b := &Bidder{
		sellerIP: sellerAddr,
	}
	b.learnAuctionRound()
	//log.Printf("DEBUG: Bidder initialized to: %v", b)
	return b
}

// Directly learn the auction round configuration from the seller along with public key
func (b *Bidder) learnAuctionRound() {
	// Get seller's public key
	uri := b.sellerIP + "/seller/key"
	response, err := http.Get(uri)
	if err != nil {
		log.Fatalf("Failed to get public key from seller: %v", err)
		os.Exit(1)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		b.sellerPublicKey = string(data)
	}

	url := b.sellerIP + "/seller/roundinfo"
	response, err = http.Get(url)
	if err != nil {
		log.Fatalf("Failed to get round information from seller: %v", err)
		os.Exit(1)
	}
	data, _ := ioutil.ReadAll(response.Body)

	var roundInfo common.AuctionRound
	err = json.Unmarshal(data, &roundInfo)
	if err != nil {
		log.Fatalf("Failed to unmarshal round information from seller: %v", err)
		os.Exit(1)
	}
	b.RoundInfo = roundInfo
}

// TODO: Get real ID of this bidder
func (b *Bidder) ProcessBid(maxBid int) {
	fakeId := big.NewInt(1999)

	var polynomials [][]*big.Int
	for price := range b.RoundInfo.Prices {
		if price <= maxBid {
			polynomials = append(polynomials, generatePolynomial(b.RoundInfo.T, fakeId, true))
		} else {
			polynomials = append(polynomials, generatePolynomial(b.RoundInfo.T, fakeId, false))
		}
	}

	fmt.Println("The following random polynomials were generated:")
	fmt.Println(polynomials)
}

// f(x) = 3x^3 + 2x + 1 => [1 2 0 3]
func generatePolynomial(degree int, id *big.Int, wantToBidThis bool) []*big.Int {
	poly := polynomial.RandomPoly(int64(degree), 5) // 5 is hard coded to make coefficients 2^5 at most

	// Change the ID
	if wantToBidThis {
		poly[0] = id
	} else {
		poly[0] = big.NewInt(0)
	}
	return poly
}
