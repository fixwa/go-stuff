package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"sync"
)

const (
	DB_HOST             = "tcp(127.0.0.1:33060)"
	DB_NAME             = "one"
	DB_USER             = "root"
	DB_PASS             = "1234"
	MAX_ROWS_PER_THREAD = 1000
	NUMBER_OF_THREADS   = 1000
)

var wg sync.WaitGroup

func main() {
	dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?charset=utf8"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	tx, transactionError := db.Begin()
	if transactionError != nil {
		panic(transactionError)
	}

	stmt, stmtError := tx.Prepare(`INSERT INTO users VALUES (?, ?, ?)`)
	if stmtError != nil {
		panic(stmtError)
	}

	wg.Add(NUMBER_OF_THREADS)
	for t := 0; t < NUMBER_OF_THREADS; t++ {
		go stmtExec(t, stmt, tx)
	}
	wg.Wait()
}

func stmtExec(threadNo int, stmt *sql.Stmt, tx *sql.Tx) {
	defer wg.Done()
	for i := 0; i < MAX_ROWS_PER_THREAD; i++ {
		stmt.Exec(uuid.NewV4().String(), randSeq(50), randSeq(255))
		//if i%10000 == 0 && i != 0 {
		fmt.Println("ThreadNumber:", threadNo, ", row:", i)
		//}
	}
	tx.Commit()
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
