package restapp

import (
	"app/internal/logic"
	"database/sql"
	"encoding/json"
)

type MySql struct {
	DB *sql.DB
}

func NewMySql(connectionString string) (*MySql, error) {
	db, err := sql.Open("mysql", connectionString)

	return &MySql{DB: db}, err
}

func (s *MySql) GetWordsFromNonIndexed(searchTerm string) (string, error) {
	return "", nil
}

func (s *MySql) GetWords(searchTerm string) (string, error) {

	searchString := searchTerm + "%"
	results, err := s.getResults(`SELECT word FROM words WHERE parsed_word like ? LIMIT 100`, searchString, 5)
	if err != nil {
		return "", err

	}
	resultsJSON, err := json.Marshal(results)
	if err != nil {
		return "", err
	}

	return string(resultsJSON), nil
}

func (s *MySql) GetWordsAsync(searchTerm string) (string, error) {
	return "", nil
}

func (s *MySql) getResults(query string, searchTerm string, maxResults int) (*[]string, error) {

	rows, err := s.DB.Query(query, searchTerm)
	if err != nil {
		return nil, err
	}
	results := []string{}
	if rows != nil {
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				return nil, err
			}
			results = append(results, name)
		}
	}
	results = *logic.Rank(searchTerm, results, maxResults)
	return &results, nil

}
