package models

import (
	"fmt"
	"sync"
	"time"
)

type LearnerApplication struct {
	Year   int
	Code   int
	Serial int
	APN    string
}

type APNGenerator struct {
	YearCounter      map[int]int
	YearCounterMutex sync.Mutex
	CodeCounterMutex sync.Mutex
	CodeCounter      int
}

func NewAPNGenerator() *APNGenerator {
	return &APNGenerator{
		YearCounter: make(map[int]int),
	}
}

func (g *APNGenerator) GenerateAPN() (*LearnerApplication, error) {
	currentYear := time.Now().Year()

	g.YearCounterMutex.Lock()
	defer g.YearCounterMutex.Unlock()

	g.YearCounter[currentYear]++
	serial := g.YearCounter[currentYear]

	g.CodeCounterMutex.Lock()
	defer g.CodeCounterMutex.Unlock()

	apn := &LearnerApplication{
		Year:   currentYear,
		Code:   3,
		Serial: serial,
		APN:    "%d%04d",
	}

	apn.APN = fmt.Sprintf(apn.APN, apn.Year, apn.Code, apn.Serial)

	return apn, nil
}
