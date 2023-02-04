package restapp

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *MySql

func errorResponse(w http.ResponseWriter, err error, errMessage string) {
	log.Fatalf(errMessage)
	log.Fatalf("Error: %v", err)

	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, "")
}

func getWords(w http.ResponseWriter, r *http.Request) {

}

//
func getWordsFromNonIndexed(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	result, err := db.GetWordsFromNonIndexed(getSearchTerm(r))
	if err != nil {
		errorResponse(w, err, "getWordsFromNonIndexed")
	}

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, result)

	elapsed := time.Since(start)
	log.Printf("Time elapsed: %v", elapsed)
}

//
func getWordsFromSql(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	result, err := db.GetWords(getSearchTerm(r))
	if err != nil {
		errorResponse(w, err, "getWords")
	}

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, result)

	elapsed := time.Since(start)
	log.Printf("Time elapsed: %v", elapsed)

}

//
func getWordsFromSqlAsync(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	result, err := db.GetWordsAsync(getSearchTerm(r))
	if err != nil {
		errorResponse(w, err, "getWordsAsync")
	}

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, result)

	elapsed := time.Since(start)
	log.Printf("Time elapsed: %v", elapsed)

}

//
func getWordsFromRedis(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	elapsed := time.Since(start)
	log.Printf("Time elapsed: %v", elapsed)

}

func getSearchTerm(r *http.Request) string {
	return r.URL.Query().Get("search")
}

func Run() {
	var err error
	db, err = NewMySql("dockeruser:dockerpass@tcp(localhost:3306)/words")
	if err != nil {
		panic(err.Error())
	}
	defer db.DB.Close()

	http.HandleFunc("/mysql/get-words", getWordsFromSql)
	http.HandleFunc("/mysql/get-words-async", getWordsFromSqlAsync)
	http.HandleFunc("/redis/get-words", getWordsFromRedis)
	http.HandleFunc("/combo/get-words", getWords)

	fmt.Printf("Listening on port 8080\n")
	http.ListenAndServe(":8080", nil)
}
