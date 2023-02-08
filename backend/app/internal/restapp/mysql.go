package restapp

import (
	"app/internal/logic"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type MySql struct {
	DB            *sql.DB
	MaxResults    int
	WordsTable    string
	SubWordsTable string
	cache         *Cache
}

type asyncResult struct {
	result []string
	err    error
}

func NewMySql(connectionString string, wordsTable string, subWordsTable string, cache *Cache) (*MySql, error) {
	db, err := sql.Open("mysql", connectionString)
	db.SetMaxOpenConns(150) //limit to 150 connections

	return &MySql{
		DB:            db,
		MaxResults:    5,
		WordsTable:    wordsTable,
		SubWordsTable: subWordsTable,
		cache:         cache,
	}, err
}

func (s *MySql) GetWords(searchTerm string) (string, error) {
	var additionalResults []string
	results, err := s.prefixSearch(s.WordsTable, searchTerm, "parsed_word")
	if err != nil {
		log.Println("Error: getStartsWith")
		return "", err
	}
	if len(results) < s.MaxResults {
		searchString := generateSearchString(searchTerm)
		statement := `SELECT word FROM %s WHERE parsed_word LIKE ? LIMIT ?`
		additionalResults, err = s.execute(s.WordsTable, statement, searchString, s.MaxResults)
		if err != nil {
			log.Println("Error: execute")
			return "", err

		}
		results = append(results, additionalResults...)
	}

	results = logic.Rank(searchTerm, results, s.MaxResults)
	return asJsonString(results)
}

func (s *MySql) GetWordsPrefixTable(searchTerm string) (string, error) {

	var additionalResults []string
	results, err := s.prefixSearch(s.WordsTable, searchTerm, "parsed_word")
	if err != nil {
		log.Println("Error: prefixSearch")
		return "", err
	}
	if len(results) < s.MaxResults {
		additionalResults, err = s.prefixSearch(s.SubWordsTable, searchTerm, "subword")
		if err != nil {
			log.Println("Error: prefixSearch")
			return "", err
		}
		results = append(results, additionalResults...)
	}
	if len(results) < s.MaxResults {
		additionalResults, err = s.fuzzySearch(s.WordsTable, searchTerm, "parsed_word")
		if err != nil {
			log.Println("Error: fuzzySearch")
			return "", err
		}
		results = append(results, additionalResults...)
	}
	if len(results) < s.MaxResults {
		additionalResults, err = s.fuzzySearch(s.SubWordsTable, searchTerm, "subword")
		if err != nil {
			log.Println("Error: fuzzySearch")
			return "", err
		}
		results = append(results, additionalResults...)
	}
	results = logic.Rank(searchTerm, results, s.MaxResults)
	return asJsonString(results)
}

func (s *MySql) GetWordsAsync(searchTerm string) (string, error) {
	var totalAsyncCalls = 10
	count, err := s.GetTableSize(s.WordsTable)
	if err != nil {
		log.Println("Error: GetTableSize")
		return "", err
	}
	searchString := generateSearchString(searchTerm)
	resultsChan := make(chan asyncResult, totalAsyncCalls)
	defer close(resultsChan)
	step := count / totalAsyncCalls
	for i := 0; i < totalAsyncCalls; i++ {
		start := i * step
		end := start + step - 1
		statement := "SELECT word FROM %s WHERE parsed_word LIKE ? AND id BETWEEN ? AND ?"
		go s.executeAsync(s.WordsTable, resultsChan, statement, searchString, start, end)
	}

	var results []string
	for i := 0; i < totalAsyncCalls; i++ {
		asyncResult := <-resultsChan
		if asyncResult.err != nil {
			log.Println("Error async call")
			return "", asyncResult.err
		}
		results = append(results, asyncResult.result...)
	}
	results = logic.Rank(searchTerm, results, s.MaxResults)
	return asJsonString(results)
}

func (s *MySql) GetTableSize(table string) (int, error) {
	tableSizeKey := fmt.Sprintf("%sTableSize", table)
	count, err := s.cache.GetInt(tableSizeKey)
	if err != nil {
		log.Println("Error: GetInt")
		return -1, err
	}

	if count == -1 || err != nil {
		query := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, s.WordsTable)
		err := s.DB.QueryRow(query).Scan(&count)
		if err != nil {
			log.Println("Error: DB Query Row")
			return -1, err
		}
		s.cache.SetInt(tableSizeKey, count, time.Hour)
	}
	return count, nil
}

////////////////////
// Helper functions
////////////////////

func (s *MySql) execute(table string, statement string, args ...interface{}) ([]string, error) {
	statement = fmt.Sprintf(statement, table)
	rows, err := s.DB.Query(statement, args...)
	if err != nil {
		log.Println("Error: DB Query")
		return nil, err
	}
	defer rows.Close()
	results := []string{}
	if rows != nil {
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				log.Println("Error: rows scan")
				return nil, err
			}
			results = append(results, name)
		}
	}
	return results, nil
}

func (s *MySql) executeAsync(table string, resultsChan chan<- asyncResult, statement string, args ...interface{}) {
	results, err := s.execute(table, statement, args...)
	if err != nil {
		log.Println("Error: execute")
		resultsChan <- asyncResult{err: err, result: nil}
		return
	}
	resultsChan <- asyncResult{err: nil, result: results}
}

func (s *MySql) prefixSearch(table string, searchTerm string, columnName string) ([]string, error) {
	searchString := searchTerm + "%"

	return s.search(table, searchString, columnName)
}

func (s *MySql) fuzzySearch(table string, searchTerm string, columnName string) ([]string, error) {
	searchString := generateFuzzySearchString(searchTerm)

	return s.search(table, searchString, columnName)
}

func (s *MySql) search(table string, searchString string, columnName string) ([]string, error) {
	start := time.Now()

	whereClause := fmt.Sprintf(`WHERE %s LIKE ? LIMIT ?`, columnName)
	statement := `SELECT word FROM %s ` + whereClause
	results, err := s.execute(table, statement, searchString, s.MaxResults)
	if err != nil {
		log.Println("Error: execute")
		return nil, err
	}
	elapsed := time.Since(start)
	log.Printf("Time elapsed for search string <%s> from table <%s> %v", searchString, table, elapsed)
	return results, nil
}

func generateSearchString(searchTerm string) string {
	searchString := ""
	for _, char := range searchTerm {
		searchString += "%" + string(char)
	}
	searchString += "%"
	return searchString
}

func generateFuzzySearchString(searchTerm string) string {
	searchString := ""
	for _, char := range searchTerm {
		searchString += string(char) + "%"
	}
	return searchString
}

func asJsonString(results []string) (string, error) {
	resultsJSON, err := json.Marshal(results)
	if err != nil {
		log.Println("Error: json Marshal")
		return "", err
	}

	return string(resultsJSON), nil
}
