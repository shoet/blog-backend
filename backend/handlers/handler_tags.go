package handlers

import (
	"net/http"

	"github.com/shoet/blog/models"
)

type TagListHandler struct {
}

func (t *TagListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tags := []*models.Tag{}
	RespondMockJSON("tags.json", &tags, w, r)
	return
}
