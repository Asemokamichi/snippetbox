package sqlite3

import (
	"database/sql"
	"time"

	"golangify.com/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (s *SnippetModel) Insert(title, content string, expires time.Time) (int, error) {
	query := `
		INSERT INTO snippets (title, content, created, expires) 
		VALUES ($1, $2, $3, $4);
	`

	result, err := s.DB.Exec(query, title, content, expires, expires.Add(time.Hour*6))
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *SnippetModel) Get(id int) (*models.Snippet, error) {
	query := `
		SELECT id, title, content, created, expires 
		FROM snippets
		WHERE expires > $1 AND id = $2
	`

	snp := &models.Snippet{}
	if err := s.DB.QueryRow(query, time.Now, id).Scan(&snp.ID, &snp.Title, &snp.Content, &snp.Expires); err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return snp, nil
}

func (s *SnippetModel) Latest() ([]*models.Snippet, error) {
	query := `
		SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > $1
		ORDER BY created DESC LIMIT 10
	`

	rows, err := s.DB.Query(query, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}
	for rows.Next() {
		sn := &models.Snippet{}
		if err := rows.Scan(&sn.ID, &sn.Title, &sn.Content, &sn.Created, &sn.Expires); err != nil {
			return nil, err
		}
		snippets = append(snippets, sn)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
