package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/internal/infrastructure/models"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/options"
	"github.com/shoet/blog/internal/usecase/create_blog"
	"github.com/shoet/blog/internal/usecase/delete_blog"
	"github.com/shoet/blog/internal/usecase/get_blog_detail"
	"github.com/shoet/blog/internal/usecase/get_blogs"
	"github.com/shoet/blog/internal/usecase/get_blogs_offset_paging"
	"github.com/shoet/blog/internal/usecase/get_tags"
	"github.com/shoet/blog/internal/usecase/put_blog"
	"github.com/shoet/blog/internal/usecase/update_public_status"
)

type BlogGetHandler struct {
	Usecase *get_blog_detail.Usecase
	jwter   JWTService
}

func NewBlogGetHandler(usecase *get_blog_detail.Usecase, jwter JWTService) *BlogGetHandler {
	return &BlogGetHandler{
		Usecase: usecase,
		jwter:   jwter,
	}
}

type BlogGetResponse struct {
	*models.Blog
}

func (l *BlogGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	blog, err := l.Usecase.Run(ctx, models.BlogId(idInt))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get blog: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}
	if blog == nil {
		response.RespondNotFound(w, r, err)
		return
	}
	// 非公開のBlogは認証が必要
	if !blog.IsPublic {
		token := r.Header.Get("Authorization")
		if token == "" {
			logger.Error("failed to get authorization header")
			response.RespondNotFound(w, r, err)
			return
		}

		if !strings.HasPrefix(token, "Bearer ") {
			logger.Error("failed to get authorization token")
			response.RespondNotFound(w, r, err)
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		_, err := l.jwter.VerifyToken(ctx, token)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to verify token: %v", err))
			response.RespondNotFound(w, r, err)
			return
		}
	}
	res := &BlogGetResponse{
		Blog: blog,
	}
	if err := response.RespondJSON(w, r, http.StatusOK, res); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type BlogAddHandler struct {
	Usecase   *create_blog.Usecase
	Validator *validator.Validate
}

func NewBlogAddHandler(
	usecase *create_blog.Usecase,
	validator *validator.Validate,
) *BlogAddHandler {
	return &BlogAddHandler{
		Usecase:   usecase,
		Validator: validator,
	}
}

func (a *BlogAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	var reqBody struct {
		Title                  string        `json:"title" validate:"required"`
		Content                string        `json:"content" validate:"required"`
		Description            string        `json:"description" validate:"required"`
		AuthorId               models.UserId `json:"authorId" validate:"required"`
		ThumbnailImageFileName string        `json:"thumbnailImageFileName"`
		IsPublic               bool          `json:"isPublic" default:"false"`
		Tags                   []string      `json:"tags" default:"[]"`
	}
	defer r.Body.Close()
	if err := response.JsonToStruct(r, &reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to parse request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	if err := a.Validator.Struct(reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	blog := &models.Blog{
		Title:                  reqBody.Title,
		Content:                reqBody.Content,
		Description:            reqBody.Description,
		AuthorId:               reqBody.AuthorId,
		ThumbnailImageFileName: reqBody.ThumbnailImageFileName,
		IsPublic:               reqBody.IsPublic,
		Tags:                   reqBody.Tags,
	}

	newBlog, err := a.Usecase.Run(ctx, blog)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to add blog: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}
	if err := response.RespondJSON(w, r, http.StatusOK, newBlog); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type BlogDeleteHandler struct {
	Usecase   *delete_blog.Usecase
	Validator *validator.Validate
}

func NewBlogDeleteHandler(
	usecase *delete_blog.Usecase,
	validator *validator.Validate,
) *BlogDeleteHandler {
	return &BlogDeleteHandler{
		Usecase:   usecase,
		Validator: validator,
	}
}

func (d *BlogDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	blogId, err := d.Usecase.Run(ctx, models.BlogId(idInt))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to delete blog: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}
	resp := struct {
		Id int `json:"id"`
	}{
		Id: int(blogId),
	}
	if err := response.RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type BlogListAdminHandler struct {
	Usecase *get_blogs.Usecase
}

func NewBlogListAdminHandler(usecase *get_blogs.Usecase) *BlogListAdminHandler {
	return &BlogListAdminHandler{
		Usecase: usecase,
	}
}

func (l *BlogListAdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)

	input := &get_blogs.GetBlogsInput{}
	resp, prevEOF, nextEOF, err := l.Usecase.Run(ctx, input)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to list blog: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}
	// TODO: 直近は管理画面ではページネーションを使わないため、EOFフラグは使わない
	_ = prevEOF
	_ = nextEOF
	if resp == nil {
		if err := response.RespondJSON(w, r, http.StatusOK, []interface{}{}); err != nil {
			logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
		}
		return
	}
	if err := response.RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type BlogListHandler struct {
	Usecase *get_blogs.Usecase
}

func NewBlogListHandler(usecase *get_blogs.Usecase) *BlogListHandler {
	return &BlogListHandler{
		Usecase: usecase,
	}
}

func (l *BlogListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)

	isPublicOnly := func() *bool { var v = true; return &v }()
	input := &get_blogs.GetBlogsInput{IsPublicOnly: isPublicOnly}
	v := r.URL.Query()
	tag := v.Get("tag")
	if tag != "" {
		input.Tag = &tag
	}
	keyword := v.Get("keyword")
	if keyword != "" {
		input.KeyWord = &keyword
	}
	cursor_id := v.Get("cursor_id") // ページネーションのカーソルID
	if cursor_id != "" {
		v, err := strconv.Atoi(cursor_id)
		if err != nil {
			err := fmt.Errorf("cursor_id is invalid")
			logger.Error(err.Error())
			response.RespondBadRequest(w, r, err)
			return
		}
		blogId := models.BlogId(v)
		input.CursorId = &blogId
	}
	direction := v.Get("direction") // ページネーションの方向
	if direction != "" {
		if direction != "prev" && direction != "next" {
			err := fmt.Errorf("direction is invalid")
			logger.Error(err.Error())
			response.RespondBadRequest(w, r, err)
			return
		}
		input.PageDirection = &direction
	}
	limit := v.Get("limit")
	if limit != "" {
		v, err := strconv.Atoi(limit)
		if err != nil {
			err := fmt.Errorf("limit is invalid")
			logger.Error(err.Error())
			response.RespondBadRequest(w, r, err)
			return
		}
		l := int64(v)
		input.Limit = &l
	}

	blogs, prevEOF, nextEOF, err := l.Usecase.Run(ctx, input)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to list blog: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}
	if blogs == nil {
		if err := response.RespondJSON(w, r, http.StatusOK, []interface{}{}); err != nil {
			logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
		}
		return
	}

	type ResponseBody struct {
		Blog    []*models.Blog `json:"blogs"`
		PrevEOF bool           `json:"prevEOF"`
		NextEOF bool           `json:"nextEOF"`
	}

	body := &ResponseBody{
		Blog:    blogs,
		PrevEOF: prevEOF,
		NextEOF: nextEOF,
	}

	if err := response.RespondJSON(w, r, http.StatusOK, body); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type BlogGetOffsetPagingHandler struct {
	Usecase *get_blogs_offset_paging.Usecase
}

func NewBlogGetOffsetPagingHandler(usecase *get_blogs_offset_paging.Usecase) *BlogGetOffsetPagingHandler {
	return &BlogGetOffsetPagingHandler{
		Usecase: usecase,
	}
}

func (l *BlogGetOffsetPagingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)

	isPublicOnly := func() *bool { var v = true; return &v }()
	input := &get_blogs_offset_paging.Input{IsPublicOnly: isPublicOnly}
	v := r.URL.Query()
	tag := v.Get("tag")
	if tag != "" {
		input.Tag = &tag
	}
	keyword := v.Get("keyword")
	if keyword != "" {
		input.KeyWord = &keyword
	}
	limit := v.Get("limit")
	if limit != "" {
		v, err := strconv.Atoi(limit)
		if err != nil {
			err := fmt.Errorf("limit is invalid")
			logger.Error(err.Error())
			response.RespondBadRequest(w, r, err)
			return
		}
		l := int64(v)
		input.Limit = &l
	}
	page := v.Get("page")
	if page != "" {
		v, err := strconv.Atoi(page)
		if err != nil {
			err := fmt.Errorf("page is invalid")
			logger.Error(err.Error())
			response.RespondBadRequest(w, r, err)
			return
		}
		p := int64(v)
		input.Page = &p
	}

	blogs, blogsCount, err := l.Usecase.Run(ctx, input)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to list blog: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}
	if blogs == nil {
		if err := response.RespondJSON(w, r, http.StatusOK, []interface{}{}); err != nil {
			logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
		}
		return
	}

	type ResponseBody struct {
		Blog       []*models.Blog `json:"blogs"`
		TotalCount int64          `json:"totalCount"`
	}

	body := &ResponseBody{
		Blog:       blogs,
		TotalCount: blogsCount,
	}

	if err := response.RespondJSON(w, r, http.StatusOK, body); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type BlogPutHandler struct {
	Usecase   *put_blog.Usecase
	Validator *validator.Validate
}

func NewBlogPutHandler(
	usecase *put_blog.Usecase,
	validator *validator.Validate,
) *BlogPutHandler {
	return &BlogPutHandler{
		Usecase:   usecase,
		Validator: validator,
	}
}

func (p *BlogPutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	var reqBody struct {
		Id                     models.BlogId `json:"id" validate:"required"`
		AuthorId               models.UserId `json:"authorId" validate:"required"`
		Title                  string        `json:"title"`
		Content                string        `json:"content"`
		Description            string        `json:"description"`
		ThumbnailImageFileName string        `json:"thumbnailImageFileName"`
		IsPublic               bool          `json:"isPublic"`
		Tags                   []string      `json:"tags"`
	}
	defer r.Body.Close()
	if err := response.JsonToStruct(r, &reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to parse request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	if err := p.Validator.Struct(reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	blog := &models.Blog{
		Id:                     reqBody.Id,
		AuthorId:               reqBody.AuthorId,
		Title:                  reqBody.Title,
		Content:                reqBody.Content,
		Description:            reqBody.Description,
		ThumbnailImageFileName: reqBody.ThumbnailImageFileName,
		IsPublic:               reqBody.IsPublic,
		Tags:                   reqBody.Tags,
	}

	newBlog, err := p.Usecase.Run(ctx, blog)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to put blog: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}
	if err := response.RespondJSON(w, r, http.StatusOK, newBlog); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type BlogUpdatePublicStatusHandler struct {
	Validator *validator.Validate
	Usecase   *update_public_status.Usecase
}

func NewBlogUpdatePublicStatusHandler(
	validator *validator.Validate,
	usecase *update_public_status.Usecase,
) *BlogUpdatePublicStatusHandler {
	return &BlogUpdatePublicStatusHandler{
		Validator: validator,
		Usecase:   usecase,
	}
}

func (h *BlogUpdatePublicStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	var reqBody struct {
		BlogId   models.BlogId `json:"blogId" validate:"required"`
		IsPublic *bool         `json:"isPublic" validate:"required"`
	}
	defer r.Body.Close()
	if err := response.JsonToStruct(r, &reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to parse request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	logger.Debug(fmt.Sprintf("reqBody: %+v", reqBody))

	if err := h.Validator.Struct(reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}
	blog, err := h.Usecase.Run(ctx, reqBody.BlogId, *reqBody.IsPublic)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to update blog public status: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}
	if blog == nil {
		response.RespondNotFound(w, r, fmt.Errorf("blog not found"))
		return
	}
	if err := response.RespondJSON(w, r, http.StatusOK, blog); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type TagListHandler struct {
	Usecase get_tags.Usecase
}

func NewTagListHandler(
	usecase get_tags.Usecase,
) *TagListHandler {
	return &TagListHandler{
		Usecase: usecase,
	}
}

func (t *TagListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	option := options.ListTagsOptions{
		Limit: 100,
	}
	resp, err := t.Usecase.Run(ctx, option)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to list tags: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}
	if resp == nil {
		if err := response.RespondJSON(w, r, http.StatusOK, []interface{}{}); err != nil {
			logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
		}
		return
	}
	if err := response.RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}
