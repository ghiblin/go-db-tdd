package godbtdd

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Repository struct {
	Db *sql.DB
}

func (r *Repository) Load(id int64) (*Blog, error) {
	query := `
		SELECT title, content, tags, created_at
		FROM blogs
		WHERE id = $1
	`
	stmt, err := r.Db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var title, content string
	var tags []string
	var createdAt time.Time
	err = stmt.QueryRow(id).Scan(&title, &content, pq.Array(&tags), &createdAt)
	if err != nil {
		return nil, err
	}
	return &Blog{
		ID:        id,
		Title:     title,
		Content:   content,
		Tags:      tags,
		CreatedAt: createdAt,
	}, nil
}

func (r *Repository) Create(blog *Blog) (int64, error) {
	query := `
		INSERT INTO blogs(title, content, tags)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	stmt, err := r.Db.Prepare(query)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	var id int64
	err = stmt.QueryRow(blog.Title, blog.Content, pq.Array(blog.Tags)).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r *Repository) Migrate() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS blogs(
		id 					serial PRIMARY KEY,
		title 			text NOT NULL,
		content 		text NOT NULL,
		tags 				text[] NOT NULL DEFAULT '{}'::text[],
		created_at 	timestamp NOT NULL DEFAULT NOW()
	)`)
	return err
}
