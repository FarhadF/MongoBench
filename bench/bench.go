package bench

import (
	//"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	//"os"
	"math/rand"
	//"sync"
	"time"
)

const (
	MongoDbHosts = "localhost:27017"
	Database     = "journaldb"
	Collection   = "journal"
)

var users = []string{"user1", "user2"}

func Bench(threads int) {
	mongoDbDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDbHosts},
		Timeout:  5 * time.Second,
		Database: Database,
	}

	mongoSession, err := mgo.DialWithInfo(mongoDbDialInfo)
	if err != nil {
		panic(err)
	}
	mongoSession.SetMode(mgo.Monotonic, true)

	ch := make(chan time.Duration)
	for query := 0; query < threads; query++ {
		go RunQuery(query, mongoSession, ch)
	}

	var x []time.Duration
	for i := 0; i < threads; i++ {

		x = append(x, <-ch)

	}

	var total time.Duration
	var n, slowest time.Duration
	for _, value := range x {
		total += value
		if value > n {
			n = value
			slowest = n
		}
	}
	fmt.Println("Average: ", total.Seconds()/float64(len(x)), "s")
	fmt.Println("Slowest: ", slowest)
}

func RunQuery(query int, mongoSession *mgo.Session, ch chan time.Duration) {
	//defer waitGroup.Done()
	sessionCopy := mongoSession.Copy()
	rand.Seed(time.Now().UnixNano())
	u := rand.Int() % len(users)

	Collection := sessionCopy.DB("journaldb").C("journal")
	defer sessionCopy.Close()
	var m bson.M
	start := time.Now()
	var q [4]bson.M
	er := bson.UnmarshalJSON([]byte(`{"branchCode": 64}`), &q[0])
	er = bson.UnmarshalJSON([]byte(`{"branchCode": 230}`), &q[1])
	er = bson.UnmarshalJSON([]byte(`{"userName":"`+users[u]+`"}`), &q[2])
	er = bson.UnmarshalJSON([]byte(`{"createDate":{$gte:ISODate("2017-05-22T05:48:15.721Z"),$lt:ISODate("2017-05-25T05:48:15.721Z")}}`), &q[3])

	if er != nil {
		panic("wtf")
	}

	n := rand.Int() % len(q)

	err := Collection.Find(q[n]).One(&m)
	dur := time.Since(start)
	if err != nil {
		log.Println("Find:", err)
	}
	fmt.Println("Query in thread ", query, " Completed. Elapsed Time: ", dur, "and query was: ", q[n])
	ch <- dur

}
