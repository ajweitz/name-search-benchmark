package restapp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var mysql *mysql

func errorResponse(w http.ResponseWriter, err error, errMessage string) {
	log.Fatalf(errMessage)
	log.Fatalf("Error: %v", err)

	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, "")
}

func getName(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	searchTerm := r.URL.Query().Get("search")

	rows, err := db.Query(`SELECT word FROM words WHERE word like ?`, searchTerm)
	if err != nil {
		errorResponse(w, err, "select query")

	}
	results := make([]string, 0, 5)
	if rows != nil {
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				errorResponse(w, err, "scanning rows")
			}
			results = append(results, name)
		}
	}

	resultsJSON, err := json.Marshal(results)
	if err != nil {
		errorResponse(w, err, "marshalling names to JSON")

	}
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(resultsJSON))

	elapsed := time.Since(start)
	log.Printf("Time elapsed: %v", elapsed)

}

func Run() {
	var err error
	mysql, err = NewMysql("dockeruser:dockerpass@tcp(localhost:3306)/words")
	if err != nil {
		panic(err.Error())
	}
	defer mysql.DB.Close()

	http.HandleFunc("/mysql/get-name", getNameFromSql)
	http.HandleFunc("/mysql/get-name-async", getNameFromSqlAsync)
	http.HandleFunc("/redis/get-name-async", getName)

	fmt.Printf("Listening on port 8080\n")
	http.ListenAndServe(":8080", nil)
}
