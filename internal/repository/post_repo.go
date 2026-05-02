package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"blog-api/internal/model"
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (r *PostRepo) Create(ctx context.Context, post *model.Post) error {
	query := `
		INSERT INTO posts (author_id, title, content, status, publish_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		post.AuthorID,
		post.Title,
		post.Content,
		post.Status,
		post.PublishAt,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create post: %w", err)
	}

	return nil
}

func (r *PostRepo) GetAllPublished(ctx context.Context) ([]model.Post, error) {
	query := `
		SELECT id, author_id, title, content, status, publish_at, created_at, updated_at
		FROM posts
		WHERE status = 'published'
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all published posts: %w", err)
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Title,
			&post.Content,
			&post.Status,
			&post.PublishAt,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan published post: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate published posts: %w", err)
	}

	return posts, nil
}

func (r *PostRepo) GetByID(ctx context.Context, id int) (*model.Post, error) {
	query := `
		SELECT id, author_id, title, content, status, publish_at, created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	var post model.Post
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.AuthorID,
		&post.Title,
		&post.Content,
		&post.Status,
		&post.PublishAt,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get post by id: %w", err)
	}

	return &post, nil
}

func (r *PostRepo) GetByAuthorID(ctx context.Context, authorID int) ([]model.Post, error) {
	query := `
		SELECT id, author_id, title, content, status, publish_at, created_at, updated_at
		FROM posts
		WHERE author_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, fmt.Errorf("get posts by author id: %w", err)
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Title,
			&post.Content,
			&post.Status,
			&post.PublishAt,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan author post: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate author posts: %w", err)
	}

	return posts, nil
}

func (r *PostRepo) Update(ctx context.Context, post *model.Post) error {
	query := `
		UPDATE posts
		SET title = $1,
		    content = $2,
		    status = $3,
		    publish_at = $4,
		    updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.Status,
		post.PublishAt,
		post.ID,
	).Scan(&post.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("update post: %w", err)
	}

	return nil
}

func (r *PostRepo) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM posts
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}

	return nil
}

func (r *PostRepo) GetScheduledToPublish(ctx context.Context, before time.Time) ([]model.Post, error) {
	query := `
		SELECT id, author_id, title, content, status, publish_at, created_at, updated_at
		FROM posts
		WHERE status = 'scheduled' AND publish_at IS NOT NULL AND publish_at <= $1
		ORDER BY publish_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, before)
	if err != nil {
		return nil, fmt.Errorf("get scheduled posts to publish: %w", err)
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Title,
			&post.Content,
			&post.Status,
			&post.PublishAt,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan scheduled post: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate scheduled posts: %w", err)
	}

	return posts, nil
}

func (r *PostRepo) PublishScheduled(ctx context.Context, id int) error {
	query := `
		UPDATE posts
		SET status = 'published',
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("publish scheduled post: %w", err)
	}

	return nil
}
