package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/puremike/social-go/internal/mailer"
	"github.com/puremike/social-go/internal/store"
)

type registerUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=72"`
}

type userWithToken struct {
	*store.UserModel
	Token string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		registerUserPayload	true	"UserModel credentials"
//	@Success		201		{object}	userWithToken		"User registered"
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

	user := &store.UserModel{
		Username: payload.Username,
		Email:    payload.Email,
	}

	ctx := r.Context()

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServer(w, r, err)
		return
	}

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	if err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.invitationExp); err != nil {

		switch err {
		case store.ErrDuplicateEmail:
			app.badRequest(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequest(w, r, err)
		default:
			app.internalServer(w, r, err)
		}
		return
	}

	// send email

	isProdEnv := app.config.environment == "production"

	vars := struct {
		Username, ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: fmt.Sprintf("%s/activate/%s", app.config.frontEndURL, plainToken),
	}

	statusCode, err := app.mailer.SendMailTrap(mailer.WelcomeUserTemplate, user.Username, user.Email, vars, !isProdEnv)

	if err != nil {
		app.logger.Errorw("Failed to send welcome email", "error", err, "statusCode", statusCode)

		//  rollback if send fails
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("Failed to rollback user after email send failure", "error", err)
		}

		app.internalServer(w, r, err)
		return
	} else {
		app.logger.Infow("Email sent", "status code", statusCode)

	}

	userWithToken := userWithToken{
		UserModel: user,
		Token:     plainToken,
	}

	if err := jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServer(w, r, err)
		return
	}
}

func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {

	type CreateEmailToken struct {
		Email    string `json:"email" validate:"required,email,max=255"`
		Password string `json:"password" validate:"required,min=6,max=72"`
	}

	var payload CreateEmailToken

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	user, err := app.store.Users.GetUserByEmail(r.Context(), payload.Email)

	if err != nil {
		app.unauthorizedError(w, r, err)
		return
	}

	if err := user.Password.Compare(payload.Password); err != nil {
		app.unauthorizedError(w, r, err)
		return
	}

	// Generate the token --> Add claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.tokenExp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.iss,
		"aud": app.config.auth.auds,
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServer(w, r, err)
	}

	if err := jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServer(w, r, err)
		return
	}

}
