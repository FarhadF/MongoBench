package bench

import (
	//"encoding/json"
	"bufio"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"os"
	//"sync"
	"time"
)

const (
	MongoDbHosts = "localhost:27017"
	Database     = "journaldb"
	Collection   = "journal"
)

func Bench(threads int, batch int, queryFilePath string) {
	mongoDbDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDbHosts},
		Timeout:  15 * time.Second,
		Database: Database,
	}

	mongoSession, err := mgo.DialWithInfo(mongoDbDialInfo)
	if err != nil {
		panic(err)
	}
	mongoSession.SetMode(mgo.Monotonic, true)

	lines, err := readLines(queryFilePath)
	/*	fmt.Println(lines)
		for _, line := range lines {
			fmt.Println("line: ")
			fmt.Println(line)
		}*/
	if err != nil {
		log.Fatal("readlines err: ", err)
	}

	b := threads / batch

	ch := make(chan time.Duration)
	var x []time.Duration

	length := len(lines)
	q := make([]bson.M, length)
	//var q []bson.M
	fmt.Println(length)
	for i := 0; i < length; i++ {
		//		fmt.Println(lines[i])
		err := bson.UnmarshalJSON([]byte(lines[i]), &q[i])

		if err != nil {
			log.Println("Error: ", err)
		}
	}

	for j := 0; j < b; j++ {
		for query := 0; query < batch; query++ {
			go RunQuery(query, j, mongoSession, ch, lines, q)
		}

		for i := 0; i < batch; i++ {

			x = append(x, <-ch)

		}
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
	fmt.Println(b)
}

func RunQuery(query int, b int, mongoSession *mgo.Session, ch chan time.Duration, lines []string, q []bson.M) {
	//defer waitGroup.Done()
	sessionCopy := mongoSession.Copy()
	rand.Seed(time.Now().UnixNano())
	//u := rand.Int() % len(users)

	Collection := sessionCopy.DB("journaldb").C("journal")
	defer sessionCopy.Close()
	var res []bson.M
	/*
		length := len(lines)
		q := make([]bson.M, length)
		for i := 0; i < length; i++ {
			fmt.Println(lines[i])
			er := bson.UnmarshalJSON([]byte(lines[i]), &q[i])

			if er != nil {
				log.Panic("Error: ", er)
			}
		}
	*/
	n := rand.Int() % len(q)
	start := time.Now()
	err := Collection.Find(q[n]).All(&res)
	dur := time.Since(start)
	if err != nil {
		log.Println("Find:", err)
	}
	fmt.Println("batch:", b, "Query in thread", query, "Completed. Elapsed Time:", dur, "and query was:", q[n])
	ch <- dur

}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
