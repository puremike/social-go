package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/puremike/social-go/internal/store"
)

func (app *application) basicAuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get the AuthHeader
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedError(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		// split the encoded header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Basic" {
			app.unauthorizedError(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		// Decode it
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		// chech the credentials
		username := app.config.auth.username
		password := app.config.auth.password

		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 || creds[0] != username || creds[1] != password {
			app.unauthorizedError(w, r, fmt.Errorf("invalid credentials"))
			return
		}

		next.ServeHTTP(w, r)
	})

}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			app.unauthorizedErrorOthers(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authToken, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedErrorOthers(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		token := parts[1]

		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedErrorOthers(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		sub, ok := claims["sub"].(float64)
		if !ok {
			app.unauthorizedErrorOthers(w, r, fmt.Errorf("invalid sub claim type"))
			return
		}
		// userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)

		userId := int(sub)

		ctx := r.Context()

		user, err := app.getUserFromCache(ctx, userId)

		if err != nil {
			app.unauthorizedErrorOthers(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, user_key, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getUserFromCache(ctx context.Context, id int) (*store.UserModel, error) {

	if !app.config.redisConfig.enabled {
		return app.store.Users.GetUserByID(ctx, id)
	}

	app.logger.Infow("cache hit", "key", "id", id)
	user, err := app.cacheStorage.Users.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		app.logger.Infow("fetching from db", "id", id)
		user, err := app.store.Users.GetUserByID(ctx, id)
		if err != nil {
			return nil, err
		}

		if err := app.cacheStorage.Users.Set(ctx, user); err != nil {
			return nil, err
		}

	}

	return user, nil
}

func (app *application) checkPostOwnership(role string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromContext(r)
		post := getPostFromContext(r)

		// check if it's the user post
		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		// role precedence check
		allowed, err := app.checkRolePrecedence(ctx, user, role)
		if err != nil {
			app.internalServer(w, r, err)
			return
		}

		if !allowed {
			app.forbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.UserModel, roleName string) (bool, error) {

	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}
