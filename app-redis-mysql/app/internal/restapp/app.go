package restapp

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func Run() {
	connectionString := "dockeruser:dockerpass@tcp(localhost:3306)/words"
	fetcher, err := NewFetchController(connectionString, "indexedwords", "words")
	if err != nil {
		log.Println("Error: NewFetchController")
		panic(err.Error())
	}
	defer fetcher.db.DB.Close()
	defer fetcher.indexedDb.DB.Close()

	http.HandleFunc("/mysql/get-words-no-index", fetcher.GetWordsFromNonIndexed) // from non-indexed table
	http.HandleFunc("/mysql/get-words", fetcher.GetWordsFromSql)                 // from indexed table
	http.HandleFunc("/mysql/get-words-async", fetcher.GetWordsFromSqlAsync)      // from indexed table, asynchronously
	http.HandleFunc("/redis/get-words", fetcher.GetWordsFromRedis)
	http.HandleFunc("/combo/get-words", fetcher.GetWords)

	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
