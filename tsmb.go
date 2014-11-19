package tsmb

import (
	"math/rand"
	"time"
)

// Every 10 minutes.
const QUERIES_PER_WEEK = 1000

const MAX_COUNT = 100
const INTERVAL = 30 * 7 * 24 * time.Hour // 1 month

type Record struct {
	T     time.Time `bson:"time"`
	Name  string    `bson:"name"`
	Count int64     `bson:"count"`
}

var Names = []string{
	"alpha",
	"beta",
	"gamma",
	"delta",
	"epsilon",
	"zeta",
	"eta",
	"theta",
	"iota",
	"kappa",
	"lambda",
	"mu",
	"nu",
	"xi",
	"omicron",
	"pi",
	"rho",
	"sigma",
	"tau",
	"upsilon",
	"phi",
	"chi",
	"psi",
	"omega",
}

func RandomName() string {
	return Names[rand.Intn(len(Names))]
}

func RandomCount() int64 {
	return int64(rand.Intn(MAX_COUNT))
}

type T struct {
	N time.Time
}

func NewT() *T {
	return &T{
		N: time.Now(),
	}
}

func (self *T) EverySecond(fn func(time.Time)) {
	secs := int64(INTERVAL.Seconds())
	for i := int64(0); i < secs; i++ {
		back := i * int64(time.Second)
		fn(self.N.Add(-time.Duration(back)))
	}
}

func (self *T) RandomRange(fn func(from, to time.Time)) {
	back := int64(INTERVAL.Nanoseconds())
	a, b := rand.Int63n(back), rand.Int63n(back)
	var min, max int64
	if a < b {
		min = a
		max = b
	} else {
		min = b
		max = a
	}
	fn(self.N.Add(-time.Duration(min)), self.N.Add(-time.Duration(max)))
}

func TimeIt(fn func()) time.Duration {
	before := time.Now()
	fn()
	after := time.Now()
	return after.Sub(before)
}
