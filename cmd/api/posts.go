package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/puremike/social-go/internal/model"
	"github.com/puremike/social-go/internal/store"
)

type postField struct {
	Content string   `json:"content" validate:"required,max=1000"`
	Title   string   `json:"title" validate:"required,max=100"`
	UserID  int      `json:"user_id"`
	Tags    []string `json:"tags"`
}

var payload postField

// createPost godoc
//
//	@Summary		Creates a post
//	@Description	Creates a post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		postField true	"Post payload"
//	@Success		201		{object}	postField
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	post := &model.PostModel{
		Content: payload.Content,
		Title:   payload.Title,
		UserID:  user.ID,
		Tags:    payload.Tags,
	}

	ctx := r.Context()
	// Create post
	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServer(w, r, err)
		return
	}

	// Response to return
	if err := jsonResponse(w, http.StatusCreated, post); err != nil {
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

	if err := jsonResponse(w, http.StatusCreated, posts); err != nil {
		app.internalServer(w, r, err)
		return
	}
}

// GetPost godoc
//
//	@Summary		Fetches a post
//	@Description	Fetches a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	postField
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [get]
func (app *application) getPostById(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)

	comments, err := app.store.Comments.GetCommentsByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServer(w, r, err)
	}

	// post.Comments = comments
	post.Comments = comments

	if err := jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServer(w, r, err)
		return
	}
}

// DeletePost godoc
//
//	@Summary		Deletes a post
//	@Description	Delete a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		204	{object} string
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [delete]
func (app *application) deletePostByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.internalServer(w, r, err)
	}

	ctx := r.Context()

	if err := app.store.Posts.DeletePostByID(ctx, id); err != nil {
		app.badRequest(w, r, err)
		return
	}

	message := "Post deleted successfully"

	if err := jsonResponse(w, http.StatusCreated, message); err != nil {
		app.internalServer(w, r, err)
		return
	}
}

// func (app *application) deleteAllPosts(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	if err := app.store.Posts.DeleteAllPosts(ctx); err != nil {
// 		app.badRequest(w, r, err)
// 		return
// 	}

// 	message := "All posts have been deleted successfully"

// 	if err := jsonResponse(w, http.StatusOK, message); err!= nil {
//         app.internalServer(w, r, err)
//         return
//     }
// }

// UpdatePost godoc
//
//	@Summary		Updates a post
//	@Description	Updates a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Post ID"
//	@Param			payload	body		UpdatePost	true	"Post payload"
//	@Success		200		{object}	model.PostModel
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [patch]

func (app *application) updatePost(w http.ResponseWriter, r *http.Request) {
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	post := &model.PostModel{
		Content: payload.Content,
		Title:   payload.Title,
		Tags:    payload.Tags,
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

	if err := jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServer(w, r, err)
		return
	}
}

type postKey string

const post_key postKey = "post"

func (app *application) postContextMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx = context.WithValue(ctx, post_key, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromContext(r *http.Request) *model.PostModel {
	post, _ := r.Context().Value(post_key).(*model.PostModel)
	return post
}
