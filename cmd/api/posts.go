package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/puremike/social-go/internal/model"
	"github.com/puremike/social-go/internal/store"
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
        app.badRequest(w, r, err)
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
		app.internalServer(w, r, err)
		return
	}

	// Response to return
	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServer(w, r, err)
		return
	}
}

func (app *application) getPostById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.internalServer(w, r, err)
        return
	}

	ctx := r.Context()

	post, err := app.store.Posts.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			app.badRequest(w, r, err)
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServer(w, r, err)
        return
	}
}