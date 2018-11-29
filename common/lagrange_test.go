package common

import (
	"log"
	"math/big"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	//make a list of compressedPoints
	var auctioneer1 CompressedPoints
	auctioneer1.Points = make(map[Price]Point)
	auctioneer1.Points[100] = Point{
		X: 1,
		Y: BigInt{
			Val: big.NewInt(1),
		},
	}
	auctioneer1.Points[200] = Point{
		X: 1,
		Y: BigInt{
			Val: big.NewInt(5),
		},
	}

	var auctioneer2 CompressedPoints
	auctioneer2.Points = make(map[Price]Point)
	auctioneer2.Points[100] = Point{
		X: 2,
		Y: BigInt{
			Val: big.NewInt(5),
		},
	}
	auctioneer2.Points[200] = Point{
		X: 2,
		Y: BigInt{
			Val: big.NewInt(7),
		},
	}

	var cps []CompressedPoints
	cps = append(cps, auctioneer1)
	cps = append(cps, auctioneer2)

	res := ComputeLagrange(cps)
	log.Println("res: ", res)

}
