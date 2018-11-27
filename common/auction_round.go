package common

import "time"

type AuctionRound struct {
	Item string
	StartTime    time.Time
	Interval     float64
	Prices       []uint
	Auctioneers  []string
	T int
	CurrentRound int
}
