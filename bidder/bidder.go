package bidder

import (
	//"log"
    _ "net/rpc"
    "net.http"
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
// Get configuration from the seller
func getConfig(sellerIP string, port int) {
    uri := sellerIP + "/seller/config"
    response, err := http.get(uri)
    if err != nil {
        log.Fatalf("Failed to get config file from seller: %v", err)
        os.Exit(1)
    }
    log.Printf("Config msg: %v", response)
}

// Query current round's price range from the seller
func (b Bidder) updatePriceRange() {
    uri := b.sellerIP + "/seller/priceRange"
    //sellerIP string
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







