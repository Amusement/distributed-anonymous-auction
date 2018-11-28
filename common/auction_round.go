package common

import "time"

type AuctionStatus uint8

const (
	BEFORE AuctionStatus = iota
	DURING
	AFTER
)

type AuctionRound struct {
	Item         string
	StartTime    time.Time
	Interval     time.Duration
	Prices       []uint
	Auctioneers  []string
	T            int
	CurrentRound int
}

func (a *AuctionRound) AuctionStatus() AuctionStatus {
	if a.afterStartTime() && !a.afterEndTime() {
		return DURING
	} else if !a.afterEndTime() {
		return BEFORE
	} else {
		return AFTER
	}
}

func (a *AuctionRound) afterEndTime() bool {
	return time.Now().UTC().After(a.StartTime.Add(a.Interval))
}

func (a *AuctionRound) afterStartTime() bool {
	return time.Now().UTC().After(a.StartTime)
}
