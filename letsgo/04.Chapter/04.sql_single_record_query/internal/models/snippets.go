package models

import (
	"database/sql"
	"errors"
	"time"
)

// The data for an individual snippet.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// Insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	stmt := `
	INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Return a specific snippet based by its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	// SQL statement
	stmt := `
		SELECT id, title, content, created, expires FROM snippets
		WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Execute SQL statement. Returns a pointer to a sql.Row object which holds
	// the result from the database.
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct.
	s := &Snippet{}

	// row.Scan() copy values from each field in sql.Row to Snippet struct.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	// Или можно записать так:
	// err := m.DB.
	// 	QueryRow(`
	// 		SELECT id, title, content, created, expires FROM snippets
	// 		WHERE expires > UTC_TIMESTAMP() AND id = ?
	// 		`).
	// 	Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		// If query returns no rows, row.Scan() will return a sql.ErrNoRows error.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord // our error
		} else {
			return nil, err
		}
	}

	return s, nil
}

// Return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
