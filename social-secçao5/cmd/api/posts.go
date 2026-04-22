package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/sikozonpc/social/internal/env/store"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	user := getUserFromContext(r)

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  user.ID,
	}

	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.InternalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	post, err := app.store.Posts.GetByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.InternalServerError(w, r, err)
		}
		return
	}

	comments, err := app.store.Comments.GetByPostID(r.Context(), id)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.InternalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)
	if post == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	if err := app.store.Posts.Delete(r.Context(), post.ID); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.InternalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	if post == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	var payload UpdatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.InternalServerError(w, r, err)
		return
	}
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")

		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.BadRequestResponse(w, r, err)
			return
		}

		post, err := app.store.Posts.GetByID(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.InternalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}
