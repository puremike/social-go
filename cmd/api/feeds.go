package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/puremike/social-go/internal/store"
)

func (app *application) getUserFeedsHandler(w http.ResponseWriter, r *http.Request) {

	// fmt.Println("Handler reached")

	// Log incoming query parameters
	// log.Println("Query Params:", r.URL.Query())

	fq := store.PagQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Tags: []string{},
		Search: "",
	}

	var err error
	fq, err = fq.Parse(r)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// Log parsed values
	// log.Printf("Parsed PagQuery: %+v\n", fq)

	// Validate after parsing
	if err := Validate.Struct(fq); err != nil {
		app.badRequest(w, r, err)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	// Log the values before calling SQL
	// log.Printf("Calling GetUserFeed with ID: %d, PagQuery: %+v\n", id, fq)

	feed, err := app.store.Posts.GetUserFeed(ctx, id, fq)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServer(w, r, err)
		return
	}
}
