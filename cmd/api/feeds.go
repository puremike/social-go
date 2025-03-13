package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) getUserFeedsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, id)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServer(w, r, err)
		return
	}
}