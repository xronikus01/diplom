package handler

import (
	"errors"
	"net/http"
	"strconv"

	"blog-api/internal/middleware"
	"blog-api/internal/service"

	"github.com/go-chi/chi/v5"
)

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

type commentRequest struct {
	Content string `json:"content"`
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	postID, err := strconv.Atoi(chi.URLParam(r, "postId"))
	if err != nil || postID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	var req commentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	comment, err := h.commentService.Create(r.Context(), userID, postID, req.Content)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			writeError(w, http.StatusBadRequest, "invalid input")
		case errors.Is(err, service.ErrPostNotFound):
			writeError(w, http.StatusNotFound, "post not found")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, comment)
}

func (h *CommentHandler) GetByPostID(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(chi.URLParam(r, "postId"))
	if err != nil || postID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	comments, err := h.commentService.GetByPostID(r.Context(), postID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			writeError(w, http.StatusBadRequest, "invalid input")
		case errors.Is(err, service.ErrPostNotFound):
			writeError(w, http.StatusNotFound, "post not found")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, comments)
}

func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	commentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || commentID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid comment id")
		return
	}

	var req commentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	comment, err := h.commentService.Update(r.Context(), userID, commentID, req.Content)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			writeError(w, http.StatusBadRequest, "invalid input")
		case errors.Is(err, service.ErrCommentNotFound):
			writeError(w, http.StatusNotFound, "comment not found")
		case errors.Is(err, service.ErrForbidden):
			writeError(w, http.StatusForbidden, "forbidden")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, comment)
}

func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	commentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || commentID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid comment id")
		return
	}

	err = h.commentService.Delete(r.Context(), userID, commentID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			writeError(w, http.StatusBadRequest, "invalid input")
		case errors.Is(err, service.ErrCommentNotFound):
			writeError(w, http.StatusNotFound, "comment not found")
		case errors.Is(err, service.ErrForbidden):
			writeError(w, http.StatusForbidden, "forbidden")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "comment deleted",
	})
}
