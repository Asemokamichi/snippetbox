package sqlite3

import (
	"database/sql"

	"golangify.com/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

// Разобарться по позже, вообще непонятно что происходит
// Глава 4.6
func (s *SnippetModel) Insert(title, content, expires string) (int, error) {
	query := `
		INSERT INTO snippets (title, content, created, expires) 
		VALUES ($1, $2, TC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL DAY = $3);
	`

	result, err := s.DB.Exec(query, title, content, expires)
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
		WHERE expires > UTC_TIMESTAMP() AND id = $1
	`

	snp := &models.Snippet{}
	if err := s.DB.QueryRow(query, id).Scan(&snp.ID, &snp.Title, &snp.Content, &snp.Expires); err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return snp, nil
}

func (s *SnippetModel) latest() ([]*models.Snippet, error) {
	return nil, nil
}
