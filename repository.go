package godbtdd

import (
	"database/sql"
	"errors"
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

func (r *Repository) fetchBlogs(query string, args ...any) ([]*Blog, error) {
	stmt, err := r.Db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blogs := make([]*Blog, 0)
	var (
		id         int64
		title      string
		content    string
		tags       []string
		created_at time.Time
	)
	for rows.Next() {
		err := rows.Scan(&id, &title, &content, pq.Array(&tags), &created_at)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, &Blog{
			ID:        id,
			Title:     title,
			Content:   content,
			Tags:      tags,
			CreatedAt: created_at,
		})
	}
	return blogs, nil
}

func (r *Repository) ListAll() ([]*Blog, error) {
	query := `
		SELECT id, title, content, tags, created_at
		FROM blogs
	`
	return r.fetchBlogs(query)
}

func (r *Repository) List(offset, limit int) ([]*Blog, error) {
	query := `
		SELECT id, title, content, tags, created_at
		FROM blogs
		LIMIT $2 OFFSET $1
	`
	return r.fetchBlogs(query, offset, limit)
}

func (r *Repository) Save(blog *Blog) error {
	// check if blog ID field is setted to its zero-value
	if blog.ID == 0 {
		// Create
		_, err := r.Create(blog)
		return err
	}

	query := `
		UPDATE blogs
		SET
			title = $2,
			content = $3,
			tags = $4
		WHERE id = $1
	`
	stmt, err := r.Db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(blog.ID, blog.Title, blog.Content, pq.Array(blog.Tags))
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Delete(id int64) error {
	return errors.New("not implemented")
}

func (r *Repository) SearchByTitle(q string, offset, limit int) ([]*Blog, error) {
	return nil, errors.New("not implemented")
}

func (r *Repository) SearchByTag(tag string, offset, limit int) ([]*Blog, error) {
	return nil, errors.New("not implemented")
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
