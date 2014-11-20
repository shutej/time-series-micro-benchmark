package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"code.google.com/p/go-uuid/uuid"
	_ "github.com/lib/pq"
	tsmb "github.com/shutej/time-series-micro-benchmark"
)

const DROP = `
DROP TABLE IF EXISTS TimeSeriesMicroBenchmark CASCADE
`

const CREATE = `CREATE TABLE TimeSeriesMicroBenchmark (
  time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  rand UUID NOT NULL,
  count BIGINT NOT NULL,
  name TEXT NOT NULL,
  PRIMARY KEY (time, rand)
)`

const INSERT = `
INSERT INTO TimeSeriesMicroBenchmark (time, rand, name, count)
     VALUES ($1, $2, $3, $4)
`

const SELECT = `
SELECT SUM(count)
  FROM TimeSeriesMicroBenchmark
 WHERE time BETWEEN $1 AND $2
   AND rand BETWEEN '00000000-0000-0000-0000-000000000000'
                AND 'ffffffff-ffff-ffff-ffff-ffffffffffff'
`

func main() {
	db, err := sql.Open("postgres", os.Getenv("POSTGRESQL_URI"))
	defer db.Close()
	if _, err := db.Exec(DROP); err != nil {
		log.Fatalf("error dropping: %v", err)
	}
	if _, err := db.Exec(CREATE); err != nil {
		log.Fatalf("error creating: %v", err)
	}

	insert, err := db.Prepare(INSERT)
	if err != nil {
		log.Fatalf("error preparing insert: %v", err)
	}

	t := tsmb.NewT()

	log.Printf("inserting events took: %v", tsmb.TimeIt(func() {
		t.EverySecond(func(t time.Time) {
			if _, err := insert.Exec(t, uuid.New(), tsmb.RandomName(), tsmb.RandomCount()); err != nil {
				log.Fatalf("error executing insert: %v", err)
			}
		})
	}))

	select_, err := db.Prepare(SELECT)
	if err != nil {
		log.Fatalf("error preparing select: %v", err)
	}

	log.Printf("querying events took: %v", tsmb.TimeIt(func() {
		for i := 0; i < tsmb.SAMPLE_SIZE; i++ {
			t.RandomRange(func(min, max time.Time) {
				row := select_.QueryRow(min, max)
				var result int64
				if err := row.Scan(&result); err != nil {
					log.Fatalf("error scanning result: %v", err)
				}
			})
		}
	}))
}
