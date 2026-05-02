package repository

import (
	"context"
	"time"

	"blog-api/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}

type PostRepository interface {
	Create(ctx context.Context, post *model.Post) error
	GetAllPublished(ctx context.Context) ([]model.Post, error)
	GetByID(ctx context.Context, id int) (*model.Post, error)
	GetByAuthorID(ctx context.Context, authorID int) ([]model.Post, error)
	Update(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, id int) error
	GetScheduledToPublish(ctx context.Context, before time.Time) ([]model.Post, error)
	PublishScheduled(ctx context.Context, id int) error
}

type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	GetByID(ctx context.Context, id int) (*model.Comment, error)
	GetByPostID(ctx context.Context, postID int) ([]model.Comment, error)
	Update(ctx context.Context, comment *model.Comment) error
	Delete(ctx context.Context, id int) error
}
