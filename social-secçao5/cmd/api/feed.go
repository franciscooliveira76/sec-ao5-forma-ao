package main

import (
	"net/http"

	"github.com/sikozonpc/social/internal/env/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	parsedFQ, err := fq.Parse(r)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(parsedFQ); err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, 341, parsedFQ)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.InternalServerError(w, r, err)
	}
}
