package main

import (
	"net/http"

	"github.com/puremike/social-go/internal/model"
)

type postField struct {
	Content string	`json:"content"`
	Title string 	`json:"title"`
	UserID int	`json:"user_id"`
	Tags []string	`json:"tags"`
}
func (app *application) CreatePost(w http.ResponseWriter, r *http.Request) {
	var payload postField

	if err := readJSON(w, r, &payload); err != nil {
        writeJSONError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

	post := &model.PostModel{
		Content : payload.Content,
		Title : payload.Title,
        UserID : payload.UserID,
        Tags : payload.Tags,
	}

	ctx := r.Context()
	// Create post
	if err := app.store.Posts.Create(ctx, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Response to return
	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}