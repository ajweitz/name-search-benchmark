package restapp

import (
	"app/internal/logic"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

type MySql struct {
	DB         *sql.DB
	MaxResults int
	Table      string
}

type asyncResult struct {
	result []string
	err    error
}

func NewMySql(connectionString string, table string) (*MySql, error) {
	db, err := sql.Open("mysql", connectionString)

	return &MySql{
		DB:         db,
		MaxResults: 5,
		Table:      table,
	}, err
}

func (s *MySql) GetWords(searchTerm string) (string, error) {
	var additionalResults []string
	results, err := s.getStartsWith(searchTerm)
	if err != nil {
		log.Println("Error: getStartsWith")
		return "", err
	}
	if len(results) < s.MaxResults {
		searchString := generateSearchString(searchTerm)
		statement := `SELECT word FROM %s WHERE length > ? AND parsed_word LIKE ? LIMIT 100`
		additionalResults, err = s.execute(statement, len(searchTerm), searchString)
		if err != nil {
			log.Println("Error: execute")
			return "", err

		}
		results = append(results, additionalResults...)
	}

	results = logic.Rank(searchTerm, results, s.MaxResults)
	return asJsonString(results)
}

func (s *MySql) GetWordsAsync(searchTerm string) (string, error) {
	var count int
	var totalAsyncCalls = 10
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, s.Table)
	err := s.DB.QueryRow(query).Scan(&count)
	if err != nil {
		log.Println("Error: DB Query Row")
		return "", err
	}
	searchString := generateSearchString(searchTerm)
	resultsChan := make(chan asyncResult, totalAsyncCalls)
	defer close(resultsChan)
	step := count / totalAsyncCalls
	for i := 0; i < totalAsyncCalls; i++ {
		start := i * step
		end := start + step - 1
		// wg.Add(1)
		statement := "SELECT word FROM %s WHERE parsed_word LIKE ? AND id BETWEEN ? AND ?"
		go s.executeAsync(resultsChan, statement, searchString, start, end)
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

////////////////////
// Helper functions
////////////////////

func (s *MySql) execute(statement string, args ...interface{}) ([]string, error) {
	statement = fmt.Sprintf(statement, s.Table)
	rows, err := s.DB.Query(statement, args...)
	if err != nil {
		log.Println("Error: DB Query")
		return nil, err
	}
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

func (s *MySql) executeAsync(resultsChan chan<- asyncResult, statement string, args ...interface{}) {
	results, err := s.execute(statement, args...)
	if err != nil {
		log.Println("Error: execute")
		resultsChan <- asyncResult{err: err, result: nil}
		return
	}
	resultsChan <- asyncResult{err: nil, result: results}
}

func (s *MySql) getStartsWith(searchTerm string) ([]string, error) {
	searchString := searchTerm + "%"
	statement := `SELECT word FROM %s WHERE parsed_word LIKE ? LIMIT ?`
	results, err := s.execute(statement, searchString, s.MaxResults)
	if err != nil {
		log.Println("Error: execute")
		return nil, err
	}
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

func asJsonString(results []string) (string, error) {
	resultsJSON, err := json.Marshal(results)
	if err != nil {
		log.Println("Error: json Marshal")
		return "", err
	}

	return string(resultsJSON), nil
}
