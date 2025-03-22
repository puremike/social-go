package main

import (
	"net/http"

	"github.com/puremike/social-go/internal/store"
)

type registerUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=72"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var payload registerUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	user := &store.UserModel {
		Username: payload.Username,
		Email: payload.Email,
	}

	ctx := r.Context()

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServer(w, r, err)
		return
	}

	if err := app.store.Users.CreateAndInvite(ctx, user, ""); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusCreated, user); err != nil {
		app.internalServer(w, r, err)
		return 
	}

	
}