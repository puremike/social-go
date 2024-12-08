package main

import (
	"net/http"

	"github.com/puremike/social-go/internal/model"
)

type userField struct {
	Username string `json:"username" validate:"required,max=12"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}


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