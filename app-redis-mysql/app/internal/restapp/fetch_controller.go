package restapp

import (
	"io"
	"log"
	"net/http"
	"time"
)

type FetchController struct {
	indexedDb *MySql
	db        *MySql
}

func NewFetchController(connectionString string, indexedTable string, nonIndexedTable string) (*FetchController, error) {
	var indexedDb *MySql
	db, err := NewMySql(connectionString, nonIndexedTable)
	if err != nil {
		return nil, err
	}
	indexedDb, err = NewMySql(connectionString, indexedTable)
	if err != nil {
		return nil, err
	}
	return &FetchController{
		indexedDb: indexedDb,
		db:        db,
	}, err
}

func (f *FetchController) GetWords(w http.ResponseWriter, r *http.Request) {

}

// Get Results from non-indexed SQL Table
func (f *FetchController) GetWordsFromNonIndexed(w http.ResponseWriter, r *http.Request) {

	f.getWordsCallback(w, r, "getWordsFromNonIndexed", f.db.GetWords)
}

// Get Results from indexed SQL Table
func (f *FetchController) GetWordsFromSql(w http.ResponseWriter, r *http.Request) {

	f.getWordsCallback(w, r, "getWordsFromSql", f.indexedDb.GetWords)

}

// Get Results from indexed SQL Table Asynchronoulsy
func (f *FetchController) GetWordsFromSqlAsync(w http.ResponseWriter, r *http.Request) {

	f.getWordsCallback(w, r, "getWordsFromSqlAsync", f.indexedDb.GetWordsAsync)
}

// Get Results by only using redis
func (f *FetchController) GetWordsFromRedis(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	elapsed := time.Since(start)
	log.Printf("Time elapsed: %v", elapsed)

}

func (f *FetchController) getWordsCallback(w http.ResponseWriter, r *http.Request, funcName string, it func(string) (string, error)) {
	start := time.Now()

	result, err := it(getSearchTerm(r))
	if err != nil {
		errorResponse(w, err, funcName)
	}

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, result)

	elapsed := time.Since(start)
	log.Printf("Time elapsed: %v", elapsed)
}

func getSearchTerm(r *http.Request) string {
	return r.URL.Query().Get("search")
}
