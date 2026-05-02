package service

import (
	"context"
	"strings"
	"time"

	"blog-api/internal/model"
	"blog-api/internal/repository"
)

type PostService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) *PostService {
	return &PostService{postRepo: postRepo}
}

func (s *PostService) Create(ctx context.Context, authorID int, title, content string, publishAt *time.Time) (*model.Post, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" || len(title) < 3 {
		return nil, ErrInvalidInput
	}

	if content == "" || len(content) < 10 {
		return nil, ErrInvalidInput
	}

	post := &model.Post{
		AuthorID:  authorID,
		Title:     title,
		Content:   content,
		PublishAt: publishAt,
	}

	if publishAt != nil && publishAt.After(time.Now()) {
		post.Status = "scheduled"
	} else {
		post.Status = "published"
	}

	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetAllPublished(ctx context.Context) ([]model.Post, error) {
	return s.postRepo.GetAllPublished(ctx)
}

func (s *PostService) GetByID(ctx context.Context, id int) (*model.Post, error) {
	if id <= 0 {
		return nil, ErrInvalidInput
	}

	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, ErrPostNotFound
	}

	return post, nil
}

func (s *PostService) GetByAuthorID(ctx context.Context, authorID int) ([]model.Post, error) {
	if authorID <= 0 {
		return nil, ErrInvalidInput
	}

	return s.postRepo.GetByAuthorID(ctx, authorID)
}

func (s *PostService) Update(ctx context.Context, userID int, postID int, title, content string, publishAt *time.Time) (*model.Post, error) {
	if postID <= 0 || userID <= 0 {
		return nil, ErrInvalidInput
	}

	existing, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrPostNotFound
	}

	if existing.AuthorID != userID {
		return nil, ErrForbidden
	}

	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" || len(title) < 3 {
		return nil, ErrInvalidInput
	}

	if content == "" || len(content) < 10 {
		return nil, ErrInvalidInput
	}

	existing.Title = title
	existing.Content = content
	existing.PublishAt = publishAt

	if publishAt != nil && publishAt.After(time.Now()) {
		existing.Status = "scheduled"
	} else {
		existing.Status = "published"
	}

	if err := s.postRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *PostService) Delete(ctx context.Context, userID int, postID int) error {
	if postID <= 0 || userID <= 0 {
		return ErrInvalidInput
	}

	existing, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPostNotFound
	}

	if existing.AuthorID != userID {
		return ErrForbidden
	}

	return s.postRepo.Delete(ctx, postID)
}

func (s *PostService) PublishScheduled(ctx context.Context, now time.Time) error {
	posts, err := s.postRepo.GetScheduledToPublish(ctx, now)
	if err != nil {
		return err
	}

	for _, post := range posts {
		if err := s.postRepo.PublishScheduled(ctx, post.ID); err != nil {
			return err
		}
	}

	return nil
}
