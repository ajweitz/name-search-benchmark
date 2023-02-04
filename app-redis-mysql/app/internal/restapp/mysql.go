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
}

func NewMySql(connectionString string) (*MySql, error) {
	db, err := sql.Open("mysql", connectionString)

	return &MySql{DB: db, MaxResults: 5}, err
}

func (s *MySql) GetWordsFromNonIndexed(searchTerm string) (string, error) {
	return s.getWordsFromTable(searchTerm, "words")

}

func (s *MySql) GetWords(searchTerm string) (string, error) {
	return s.getWordsFromTable(searchTerm, "indexedwords")
}

func (s *MySql) getWordsFromTable(searchTerm string, table string) (string, error) {

	var additionalResults []string
	searchString := searchTerm + "%"
	statement := `SELECT word FROM %s WHERE parsed_word LIKE ? LIMIT ?`
	statement = fmt.Sprintf(statement, table)
	results, err := s.getResults(searchTerm,
		s.MaxResults,
		false,
		statement,
		searchString, s.MaxResults)
	if err != nil {
		log.Println("Error: getResults")
		return "", err
	}
	if len(results) < s.MaxResults {
		searchString = ""
		for _, char := range searchTerm {
			searchString += "%" + string(char)
		}
		searchString += "%"
		statement = `SELECT word FROM %s WHERE length > ? AND parsed_word LIKE ? LIMIT 100`
		statement = fmt.Sprintf(statement, table)
		additionalResults, err = s.getResults(searchTerm,
			s.MaxResults,
			true,
			statement,
			len(searchTerm), searchString)
		if err != nil {
			log.Println("Error: getResults 2nd time")
			return "", err

		}
		additionalResults = additionalResults[0:logic.Min(len(additionalResults), s.MaxResults-len(results))]
		results = append(results, additionalResults...)

	}

	resultsJSON, err := json.Marshal(results)
	if err != nil {
		log.Println("Error: json Marshal")
		return "", err
	}

	return string(resultsJSON), nil
}
func (s *MySql) GetWordsAsync(searchTerm string) (string, error) {
	return "", nil
}

func (s *MySql) getResults(searchTerm string, maxResults int, rank bool, statement string, args ...interface{}) ([]string, error) {

	rows, err := s.DB.Query(statement, args...)
	if err != nil {
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

	if rank {
		results = *logic.Rank(searchTerm, results, maxResults)
	}
	results = results[0:logic.Min(len(results), s.MaxResults)]
	return results, nil

}
