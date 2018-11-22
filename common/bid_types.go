package common

import (
    "math/big"
)

// Data is big.Int, since our IDs can get huge
type Point struct {
	X big.Int
	Y big.Int
}

type BidPoints struct {
	BidderID string
	Points   []Point
}
