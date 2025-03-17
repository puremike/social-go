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

type userField struct {
	Username string `json:"username" validate:"required,max=12"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type FollowUser struct {
	UserID int `json:"user_id"`
}


type userKey string
const user_key userKey = "user"

// GetUser godoc
//	@Summary		Fetch a user profile
//	@Description	get user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	model.UserModel
//	@Failure		400	{object}	httputil.HTTPError
//	@Failure		404	{object}	httputil.HTTPError
//	@Failure		500	{object}	httputil.HTTPError
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var payload userField
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
        return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
        return
	}

	user := &model.UserModel{
		Username: payload.Username,
        Email:    payload.Email,
        Password: payload.Password,
	}

	ctx := r.Context()

	if err := app.store.Users.Create(ctx, user); err != nil {
		app.internalServer(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusCreated, user); err != nil {
		app.internalServer(w, r, err)
		return 
	}
}

func (app *application) getUserByID(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServer(w, r, err)
		return
	}

}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{id}/follow [put]

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerID := getUserFromContext(r)

	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Follow(ctx, followerID.ID, payload.UserID); err != nil {
		app.internalServer(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServer(w, r, err)
		return
	}
}

// UnfollowUser gdoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{id}/unfollow [put]

func (app *application) unFollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unFollowerID := getUserFromContext(r)
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Unfollow(ctx, unFollowerID.ID, payload.UserID); err != nil {
		app.internalServer(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServer(w, r, err)
		return
	}
}

func (app *application) userContextMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			app.internalServer(w, r, err)
		return
		}
		ctx := r.Context()
		user, err := app.store.Users.GetUserByID(ctx, id)
		if err != nil {
			if errors.Is(err, store.ErrUserNotFound) {
				app.notFound(w, r, err)
				return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		}

		ctx = context.WithValue(ctx, user_key, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *model.UserModel {
	user, _ := r.Context().Value(user_key).(*model.UserModel)
	return user
}