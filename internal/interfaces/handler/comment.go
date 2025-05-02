package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/shoet/blog/internal/infrastructure/models"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/usecase/get_comments"
	"github.com/shoet/blog/internal/usecase/post_comment"
)

type GetCommentsHandler struct {
	Usecase *get_comments.Usecase
}

func NewGetCommentsHandler(usecase *get_comments.Usecase) *GetCommentsHandler {
	return &GetCommentsHandler{
		Usecase: usecase,
	}
}

type GetCommentsResponse struct {
	Comments []*models.Comment `json:"comments"`
}

/*
RequestBody:

	path: /blogs/{id}/comments

Response:

	comments: []Comment
		commentId: int
		blogId: int
		clientId: string | null
		userId: int | null
		content: string
		isEdited: bool
		isDeleted: bool
		threadId: string | null
		created: time.Time
		modified: time.Time
*/
func (h *GetCommentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	id := chi.URLParam(r, "id")
	if id == "" {
		logger.Error("failed to get id from url")
		response.RespondBadRequest(w, r, nil)
		return
	}
	idInt, err := strconv.Atoi(strings.TrimSpace(id))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to convert id to int: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}
	comments, err := h.Usecase.Run(ctx, models.BlogId(idInt))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get comments: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}
	res := GetCommentsResponse{
		Comments: comments,
	}
	if err := response.RespondJSON(w, r, http.StatusOK, res); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type PostCommentHandler struct {
	Usecase   *post_comment.Usecase
	jwter     JWTService
	Validator *validator.Validate
}

func NewPostCommentHandler(usecase *post_comment.Usecase, jwter JWTService, validator *validator.Validate) *PostCommentHandler {
	return &PostCommentHandler{
		Usecase:   usecase,
		jwter:     jwter,
		Validator: validator,
	}
}

type PostCommentRequest struct {
	UserId          *models.UserId `json:"userId"`
	ClientId        *string        `json:"clientId"`
	Content         string         `json:"content" validate:"required"`
	ThreadCommentId *int64         `json:"threadCommentId,omitempty"`
}

type PostCommentResponse struct {
	CommentId models.CommentId `json:"commentId"`
}

/*
RequestBody:

	path: /blogs/{id}/comments

	application/json:
		userId: int | null
		clientId: string | null
		content: string
		threadCommentId: int | null

Response:

	commentId: int
*/
func (h *PostCommentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	id := chi.URLParam(r, "id")
	if id == "" {
		logger.Error("failed to get id from url")
		response.RespondBadRequest(w, r, nil)
		return
	}
	idInt, err := strconv.Atoi(strings.TrimSpace(id))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to convert id to int: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	defer r.Body.Close()

	var req PostCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(fmt.Sprintf("failed to decode request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}
	if err := h.Validator.Struct(req); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}
	if req.UserId == nil && req.ClientId == nil {
		logger.Error("client_id or user_id is required")
		response.RespondBadRequest(w, r, nil)
		return
	}

	// UserIdによる投稿は認証が必要
	if req.UserId != nil {
		token := r.Header.Get("Authorization")
		if token == "" {
			logger.Error("failed to get authorization header")
			response.RespondUnauthorized(w, r, err)
			return
		}

		if !strings.HasPrefix(token, "Bearer ") {
			logger.Error("failed to get authorization token")
			response.RespondUnauthorized(w, r, err)
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		userId, err := h.jwter.VerifyToken(ctx, token)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to verify token: %v", err))
			response.RespondUnauthorized(w, r, err)
			return
		}
		if userId != *req.UserId {
			logger.Error("user_id is not matched")
			response.RespondUnauthorized(w, r, err)
			return
		}
	}

	commentId, err := h.Usecase.Run(ctx, models.BlogId(idInt), req.UserId, req.ClientId, req.ThreadCommentId, req.Content)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to post comment: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}
	res := PostCommentResponse{
		CommentId: commentId,
	}
	if err := response.RespondJSON(w, r, http.StatusOK, res); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}
