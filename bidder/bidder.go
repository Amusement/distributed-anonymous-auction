package bidder

import (
	//"log"
    _ "net/rpc"
)


type Bidder struct {
    // List of clients
    listOfAuctioneers string
    publicKey string
    sellerIP string
    startTime string
}

func InitBidder(sellerAddr string) Bidder {
    //var rawConfig string

    // unmarshal rawConfig, return a bidder
    return Bidder{
        //listOfAuctioneers: ["TODO"]
        //publicKey: "TODO",
        sellerIP: sellerAddr,
    }
}


// =============== REST call to seller ===============
func getConfig() {
}

func (b Bidder) updatePriceRange() {
}

func (b Bidder) getRoundStatus() {

}


// =============== RPC call to auctioneer ===========

func (b Bidder) sendPolynomialPoint() {
    // loop through all auctioneers and send points
}

// =============== bidder functions ================
func (b Bidder) generateID(auctioneerID, price int) {
    // return big.Int from math/big library?
}


// ============== Other functions ==================
func generatePolynomail() {

}







