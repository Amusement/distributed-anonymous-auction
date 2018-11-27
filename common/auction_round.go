package common

import "time"

type AuctionRound struct {
	StartTime    time.Time
	Interval     time.Duration
	Prices       []uint
	Auctioneers  []string
	CurrentRound int
}
