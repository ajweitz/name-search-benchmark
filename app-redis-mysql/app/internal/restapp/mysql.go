package restapp

import (
	"app/internal/logic"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type MySql struct {
	DB         *sql.DB
	MaxResults int
	Table      string
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
		// additionalResults = additionalResults[0:logic.Min(len(additionalResults), s.MaxResults-len(results))]
		results = append(results, additionalResults...)
	}

	results = logic.Rank(searchTerm, results, s.MaxResults)

	// results = results[0:logic.Min(len(results), s.MaxResults)]

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
	var wg sync.WaitGroup
	resultsChan := make(chan []string, 10)
	step := count / totalAsyncCalls
	for i := 0; i < 10; i++ {
		start := i * step
		end := start + step - 1
		wg.Add(1)
		statement := "SELECT name FROM %s WHERE name LIKE ? AND id BETWEEN ? AND ?"

		go s.executeAsync(&wg, resultsChan, statement, start, end)
	}
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var results []string

	for r := range resultsChan {
		results = append(results, r...)
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

func (s *MySql) executeAsync(wg *sync.WaitGroup, resultsChan chan<- []string, statement string, args ...interface{}) {
	defer wg.Done()
	results, err := s.execute(statement, args...)
	if err != nil {
		log.Println("Error: execute")
		return
	}
	resultsChan <- results
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
