// Package extractors contains extractor functions for getting data out
// of the request.
package extractors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hay-kot/hookfeed/backend/internal/core/apperrors"
	"github.com/hay-kot/hookfeed/backend/internal/core/validate"
)

// Body decodes the request body into a struct and validates the struct using
// the `validate` tags
func Body[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		if errors.Is(err, io.EOF) {
			return v, apperrors.ErrNoBody
		}

		return v, err
	}

	err := validate.Check(v)
	if err != nil {
		return v, err
	}

	return v, err
}

// String extracts the key field from the path and returns an error if the value is empty
func String(r *http.Request, key string) (string, error) {
	v := chi.URLParam(r, key)

	if v == "" {
		return "", validate.NewRouteKeyErrorWithMessage(key, "no value provided")
	}

	return v, nil
}

// String2 extracts the keys field from the path and returns an error if the value is empty
func String2(r *http.Request, key1, key2 string) (string, string, error) {
	v1 := chi.URLParam(r, key1)
	if v1 == "" {
		return "", "", validate.NewRouteKeyErrorWithMessage(key1, "no value provided")
	}

	v2 := chi.URLParam(r, key2)
	if v2 == "" {
		return "", "", validate.NewRouteKeyErrorWithMessage(key2, "no value provided")
	}

	return v1, v2, nil
}

// ID extracts the id field from the path, validates it, and returns it as
// a ID. If the id is not a valid ID, an error is returned.
func ID(r *http.Request, key string) (uuid.UUID, error) {
	v := chi.URLParam(r, key)

	id, err := uuid.Parse(v)
	if err != nil {
		return uuid.Nil, validate.NewRouteKeyErrorWithMessage(key, "unable to parse UUID")
	}

	return id, nil
}

// ID2 extracts two IDs from the path and validates them. If either of the IDs
// is not a valid ID, an error is returned.
func ID2(r *http.Request, key1, key2 string) (uuid.UUID, uuid.UUID, error) {
	id1, err := ID(r, key1)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	id2, err := ID(r, key2)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return id1, id2, nil
}

// ID3 extracts two IDs from the path and validates them. If either of the IDs
// is not a valid ID, an error is returned.
func ID3(r *http.Request, key1, key2, key3 string) (uuid.UUID, uuid.UUID, uuid.UUID, error) {
	id1, err := ID(r, key1)
	if err != nil {
		return uuid.Nil, uuid.Nil, uuid.Nil, err
	}

	id2, err := ID(r, key2)
	if err != nil {
		return uuid.Nil, uuid.Nil, uuid.Nil, err
	}

	id3, err := ID(r, key3)
	if err != nil {
		return uuid.Nil, uuid.Nil, uuid.Nil, err
	}

	return id1, id2, id3, nil
}

func Slug(r *http.Request, key string) (string, error) {
	slug := chi.URLParam(r, key)
	if slug == "" {
		return "", validate.NewRouteKeyError(key)
	}
	return slug, nil
}

// BodyWithID combines the calls of Body and ID into one call and extracts both
// the ID and the body of the request.
func BodyWithID[T any](r *http.Request, key string) (uuid.UUID, T, error) {
	id, err := ID(r, key)
	if err != nil {
		var v T
		return uuid.Nil, v, err
	}

	body, err := Body[T](r)
	if err != nil {
		var v T
		return uuid.Nil, v, err
	}

	return id, body, nil
}

// IntID extracts the id field from the path, validates it, and returns it as
// a ID. If the id is not a valid ID, an error is returned.
func IntID(r *http.Request, key string) (int, error) {
	v := chi.URLParam(r, key)

	id, err := strconv.Atoi(v)
	if err != nil {
		// Handle error
		fmt.Println(err)
		return 0, validate.NewRouteKeyErrorWithMessage(key, "unable to parse int")
	}

	return id, nil
}
