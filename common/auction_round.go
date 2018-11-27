package common

import "time"

type AuctionRound struct {
	StartTime    time.Time
	Interval     time.Duration
	Prices       []int
	Auctioneers  []string
	CurrentRound int
}
