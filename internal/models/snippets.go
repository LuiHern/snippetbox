package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets(title, content, created, expires)
					 VALUES($1, $2, current_timestamp, current_timestamp + make_interval(days := $3)) RETURNING id`

	row := m.DB.QueryRow(stmt, title, content, expires)
	var id int

	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created, expires 
				   FROM snippets 
				   WHERE expires > CURRENT_TIMESTAMP
	         AND id = $1`

	// query for one row, errors are handled when scan() is called
	row := m.DB.QueryRow(stmt, id)
	// init zeroed snippet,
	var s Snippet

	// map row to Snippet
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// check if query returned no rows or other error
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires
					 FROM snippets
					 WHERE expires > CURRENT_TIMESTAMP
					 ORDER BY id DESC LIMIT 10`

	// Query returns multiple rows
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// close the result set from Query, should come after you check for an error, panics if result is nil
	defer rows.Close()
	var snippets []Snippet

	for rows.Next() {
		var s Snippet
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
