package main

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var dbConn *sql.DB

var wg sync.WaitGroup

type QueueData struct {
	id      int
	message string
	status  string
}

func init() {
	// initiate DB conn
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/sqsbroker")
	if err != nil {
		log.Fatal(err)
	}
	dbConn = db
}

// func insertInQueue(message string) {
// 	rand.Int()
// 	_, err := dbConn.Exec("INSERT INTO sbroker (message,status) VALUES (?,'todo')", message)
// 	if err != nil {
// 		log.Printf("Something went wrong while inserting message in DB : %v", err.Error())
// 	}
// }

func consumer() {
	// start transaction
	txn, txnerr := dbConn.Begin()
	if txnerr != nil {
		log.Printf("Error starting txn : %v", txnerr.Error())
	}

	// select all message those having status as todo. Lock the row to avoid dual updates
	row := txn.QueryRow("SELECT * FROM sbroker WHERE status='todo' ORDER BY id LIMIT 1 FOR UPDATE SKIP LOCKED")
	if row.Err() != nil {
		log.Printf("Error while executing select query : %v", row.Err())
	}

	var qData QueueData

	scanerr := row.Scan(&qData.id, &qData.message, &qData.status)
	if scanerr != nil {
		log.Printf("ERROR Scan : %v", scanerr.Error())
	}

	// consume and print message
	log.Printf("Message : %s", qData.message)

	// update status to done for consumed message
	_, updateerr := txn.Exec("UPDATE sbroker SET status='done' WHERE id=?", qData.id)
	if updateerr != nil {
		log.Printf("Error while updating status : %v", updateerr.Error())
	}

	// commit transaction
	commiterr := txn.Commit()
	if commiterr != nil {
		log.Printf("Error while commiting transaction : %v", commiterr.Error())
	}

	wg.Done()
}

func main() {
	log.Println("SQS borker using MySQL")

	// insert into queue
	// insertInQueue("B+ tree")
	wg.Add(5)

	// run consumers
	for i := 0; i < 5; i++ {
		go consumer()
	}

	// wait for consumers to finish the tasks
	wg.Wait()
}
