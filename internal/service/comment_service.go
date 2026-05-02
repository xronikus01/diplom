package service

import (
	"context"
	"strings"

	"blog-api/internal/model"
	"blog-api/internal/repository"
)

type CommentService struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
}

func NewCommentService(commentRepo repository.CommentRepository, postRepo repository.PostRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

func (s *CommentService) Create(ctx context.Context, userID int, postID int, content string) (*model.Comment, error) {
	if userID <= 0 || postID <= 0 {
		return nil, ErrInvalidInput
	}

	content = strings.TrimSpace(content)
	if len(content) < 2 {
		return nil, ErrInvalidInput
	}

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, ErrPostNotFound
	}

	comment := &model.Comment{
		PostID:   postID,
		AuthorID: userID,
		Content:  content,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) GetByPostID(ctx context.Context, postID int) ([]model.Comment, error) {
	if postID <= 0 {
		return nil, ErrInvalidInput
	}

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, ErrPostNotFound
	}

	return s.commentRepo.GetByPostID(ctx, postID)
}

func (s *CommentService) Update(ctx context.Context, userID int, commentID int, content string) (*model.Comment, error) {
	if userID <= 0 || commentID <= 0 {
		return nil, ErrInvalidInput
	}

	content = strings.TrimSpace(content)
	if len(content) < 2 {
		return nil, ErrInvalidInput
	}

	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, ErrCommentNotFound
	}

	if comment.AuthorID != userID {
		return nil, ErrForbidden
	}

	comment.Content = content

	if err := s.commentRepo.Update(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) Delete(ctx context.Context, userID int, commentID int) error {
	if userID <= 0 || commentID <= 0 {
		return ErrInvalidInput
	}

	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return ErrCommentNotFound
	}

	if comment.AuthorID != userID {
		return ErrForbidden
	}

	return s.commentRepo.Delete(ctx, commentID)
}
