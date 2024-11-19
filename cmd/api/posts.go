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
	Content string	`json:"content" validate:"required,max=1000"`
	Title string 	`json:"title" validate:"required,max=100"`
	UserID int	`json:"user_id"`
	Tags []string	`json:"tags"`
}
	
func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	payload := postField{}

	if err := readJSON(w, r, &payload); err != nil {
        app.badRequest(w, r, err)
        return
    }

	if err := Validate.Struct(payload); err != nil {
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

func (app *application) getAllPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	posts, err := app.store.Posts.GetAllPosts(ctx)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, posts); err != nil {
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
			app.notFound(w, r, err)
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}

	comments, err := app.store.Comments.GetCommentsByPostID(ctx, id)
	if err != nil {
		app.internalServer(w, r, err)
	}

	// post.Comments = comments
	post.Comments = comments

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServer(w, r, err)
        return
	}
}

func (app *application) deletePostByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.internalServer(w, r, err)
	}

	ctx := r.Context()

	message, err := app.store.Posts.DeletePostByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			app.notFound(w, r, err)
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}

	if err := writeJSON(w, http.StatusCreated, message); err != nil {
		app.internalServer(w, r, err)
		return
	}
}

func (app *application) deleteAllPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	message, err := app.store.Posts.DeleteAllPosts(ctx)

	if err != nil {
		app.badRequest(w, r, err)
	}
	
	if err := writeJSON(w, http.StatusOK, message); err!= nil {
        app.internalServer(w, r, err)
        return
    }
}

func (app *application) updatePost(w http.ResponseWriter, r *http.Request) {

	payload := postField{}


	if err := readJSON(w, r, &payload); err != nil {
        app.badRequest(w, r, err)
        return
    }

	if err := Validate.Struct(payload); err!= nil {
        app.badRequest(w, r, err)
        return
    }
	
	post := &model.PostModel{
		Content: payload.Content,
        Title: payload.Title,
        Tags: payload.Tags,
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.internalServer(w, r, err)
	}

	ctx := r.Context()

	if err = app.store.Posts.UpdatePost(ctx, id, post); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServer(w, r, err)
		return
	}
}