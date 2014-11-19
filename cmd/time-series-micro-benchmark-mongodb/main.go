package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	tsmb "github.com/shutej/time-series-micro-benchmark"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getSession() (*mgo.Session, error) {
	env := os.Getenv("MONGODB_URI")
	if env == "" {
		return nil, fmt.Errorf("MONGODB_URI is missing")
	}
	log.Printf("connecting to mongodb database: %s", env)
	url, err := url.Parse(env)
	if err != nil {
		return nil, err
	}

	database := ""
	if url.Path != "" {
		database = url.Path[1:len(url.Path)]
	}

	username := ""
	password := ""
	if url.User != nil {
		username = url.User.Username()
		password, _ = url.User.Password()
	}

	info := &mgo.DialInfo{
		Addrs:    []string{url.Host},
		Username: username,
		Password: password,
		Database: database,
	}
	return mgo.DialWithInfo(info)
}

func main() {
	session, err := getSession()
	if err != nil {
		log.Fatal(err)
	}

	db := session.DB("time-series-micro-benchmark")
	c := db.C("events")
	c.DropCollection()

	t := tsmb.NewT()

	log.Printf("inserting events took: %v", tsmb.TimeIt(func() {
		t.EverySecond(func(t time.Time) {
			c.Insert(tsmb.Record{
				T:     t,
				Count: tsmb.RandomCount(),
				Name:  tsmb.RandomName(),
			})
		})
	}))

	index := mgo.Index{
		Key:      []string{"time"},
		Unique:   true,
		DropDups: true,
	}
	if err = c.EnsureIndex(index); err != nil {
		log.Fatal(err)
	}

	// Ensure your design works if scale changes by 10X or 20X but the right
	// solution for X often not optimal for 100X.
	// http://static.googleusercontent.com/media/research.google.com/en/us/people/jeff/stanford-295-talk.pdf
	for i := 0; i < 20*tsmb.QUERIES_PER_WEEK; i++ {
		t.RandomRange(func(min, max time.Time) {
			pipe := c.Pipe([]bson.M{
				bson.M{"$match": bson.M{"time": bson.M{"$gt": min, "$lt": max}}},
				bson.M{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": "$count"}}},
			})
			var result struct {
				Count int64 `bson:"count"`
			}
			if err = pipe.One(&result); err != nil {
				log.Fatal(err)
			}
			log.Printf("count: %v", result.Count)
		})
	}
}
