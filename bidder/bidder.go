package bidder

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
	"crypto/rsa"
	//"bytes"
	"bytes"
)

type Bidder struct {
	RoundInfo		common.AuctionRound			// Initially retrieved round info.

	secretID        int
	sellerPublicKey *rsa.PublicKey
	sellerIP        string
	bidderIP string
}

func InitBidder(sellerAddr string, bidderIP string) *Bidder {
	b := &Bidder{
		sellerIP: sellerAddr,
		bidderIP: bidderIP,
	}
	b.learnAuctionRound()
	//log.Printf("Bidder initialized to: %v", b)
	return b
}

// Directly learn the auction round configuration from the seller along with public key
func (b *Bidder) learnAuctionRound() {
	// Get seller's public key
	uri := b.sellerIP + "/seller/key"
	response, err := http.Get(uri)
	if err != nil {
		log.Fatalf("Failed to get public key from seller: %v", err)
	}

	// Parse and store seller's public key
	data, _ := ioutil.ReadAll(response.Body)
	key, err := common.UnmarshalPemToKey(data)
	if err != nil {
		log.Fatalf(err.Error())
	}
	b.sellerPublicKey = key

	// Get auction round info
	url := b.sellerIP + "/seller/roundinfo"
	response, err = http.Get(url)
	if err != nil {
		log.Fatalf("Failed to get round information from seller: %v", err)
	}

	// Parse and store auction round info
	data, _ = ioutil.ReadAll(response.Body)
	var roundInfo common.AuctionRound
	err = json.Unmarshal(data, &roundInfo)
	if err != nil {
		log.Fatalf("Failed to unmarshal round information from seller: %v", err)
		os.Exit(1)
	}
	b.RoundInfo = roundInfo
}

// TODO: Choosing a hardcoded port for now, not ideal
func (b *Bidder) ProcessBid(maxBid int) {
	maxBidU := uint(maxBid)

	// Choose a port, TODO: make port part of Bidder struct? Accessed frequently
	chosenPort := 4331
	fmt.Println("Chose port: ", chosenPort)

	var polynomials []polynomial.Poly
	for _, price := range b.RoundInfo.Prices {
		if price <= maxBidU {
			id := b.selfIdentify(chosenPort, price)
			polynomials = append(polynomials, generatePolynomial(b.RoundInfo.T, id))
		} else {
			polynomials = append(polynomials, generatePolynomial(b.RoundInfo.T, nil))
		}
	}

	//fmt.Println("The following random polynomials were generated:\n", polynomials)
	b.samplePoints(polynomials)
}

// Prepare, but do not send, the points for each auctioneer
// Important note: auctioneers are given points corresponding to their order
func (b *Bidder) samplePoints(polynomials []polynomial.Poly) {
	// example entry: 1:500:[(1,2)]
	//				  1:700:[(1,19)]
	auctioneerPricePoints := make(map[int]map[common.Price]common.Point)

	for i, _ := range b.RoundInfo.Auctioneers {
		x := i+1
		// Initialize nested map
		auctioneerPricePoints[x] = make(map[common.Price]common.Point)
		for _, price := range b.RoundInfo.Prices {
			// Evaluate polynomial for this price at the point x
			bigX := big.NewInt(int64(x))
			y := common.BigInt{
				polynomials[i].Eval(bigX, nil),
			}

			sampledPoint := common.Point{
				x,
				y,
			}

			auctioneerPricePoints[i+1][common.Price(price)] = sampledPoint
		}
	}

	//fmt.Printf("The following points were sampled:\n%v", auctioneerPricePoints)
	b.sendPoints(auctioneerPricePoints)
}

func (b *Bidder) sendPoints(auctioneerPricePoints map[int]map[common.Price]common.Point) {
	failed := false

	for i, auctioneer := range b.RoundInfo.Auctioneers {
		bidPoints := common.BidPoints{
			BidderID: b.bidderIP,			// TODO
			Points: auctioneerPricePoints[i],
		}

		// Internal test
		bidPointsEnc, err := common.MarshalBidPoints(bidPoints)
		if err != nil {
			fmt.Println("Error encoding bidPoints: ", err)
			failed = true
			break
		}

		var bidPointsDec common.BidPoints
		err = common.UnmarshalBidPoints(bidPointsEnc, &bidPointsDec)
		if err != nil {
			fmt.Println("Error decoding bidPoints: ", err)
			failed = true
			break
		}

		// Rough check of equality
		if bidPointsDec.Points[0] != bidPoints.Points[0] {
			fmt.Println("Didn't get the same result after decoding encoded points.")
			failed = true
			break
		}

		url := "http://" + auctioneer + "/auctioneer/sendBid"
		//req, err := http.NewRequest("POST", url, bytes.NewBuffer(bidPointsEnc))
		client := http.DefaultClient
		//resp, err := client.Do(req)
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(bidPointsEnc))
		if err != nil {
			fmt.Printf("Unable to reach auctioneer %v\n", auctioneer)
			failed = true
			break
		}

		if resp.StatusCode != 200 {
			fmt.Printf("Auctioneer %v rejected the bid.\n", auctioneer)
			failed = true
			break
		}
	}

	if failed {
		fmt.Println("One or more auctioneers was unreachable or rejected the bid.")
	} else {
		fmt.Println("All auctioneers accepted the bid.")
	}
}

// Make a secretID using the given port for given price
func (b *Bidder) selfIdentify(chosenPort int, price uint) *big.Int {
	localIPPort := fmt.Sprintf("%v:%v", b.bidderIP, chosenPort)

	encryptedIDBytes, err := common.EncryptID(localIPPort, price, b.sellerPublicKey)
	if err != nil {
		// TODO: Error handling?
		log.Fatalf("Failed to encrypt bidder secretID: %v", err)
	}

	id := big.NewInt(0)
	id = id.SetBytes(encryptedIDBytes)

	return id
}

// f(x) = 3x^3 + 2x + 1 => [1 2 0 3]
func generatePolynomial(degree int, id *big.Int) polynomial.Poly {
	poly := polynomial.RandomPoly(int64(degree), 5) // 5 is hard coded to make coefficients 2^5 at most

	// Change the ID
	if id != nil {
		poly[0] = id
	} else {
		poly[0] = big.NewInt(0)
	}
	return poly
}
