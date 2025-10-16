package extractors

import (
	"net/http"
	"reflect"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/hay-kot/hookfeed/backend/internal/core/validate"
)

var queryDecoder = schema.NewDecoder()

// init is required for queryDecoder and there are no side effects
func init() { //nolint:gochecknoinits
	queryDecoder.IgnoreUnknownKeys(true)

	queryDecoder.RegisterConverter(uuid.UUID{}, func(s string) reflect.Value {
		v, err := uuid.Parse(s)
		if err != nil {
			// TODO: what to do here?
			v = uuid.Nil
		}
		return reflect.ValueOf(v)
	})

	queryDecoder.RegisterConverter(civil.Date{}, func(s string) reflect.Value {
		v, err := civil.ParseDate(s)
		if err != nil {
			// TODO: what to do here?
			v = civil.Date{}
		}
		return reflect.ValueOf(v)
	})
}

func QueryT[T any](r *http.Request) (T, error) {
	var v T
	err := Query(r, &v)
	return v, err
}

func Query(r *http.Request, v any) error {
	err := queryDecoder.Decode(v, r.URL.Query())
	if err != nil {
		return err
	}

	return validate.Check(v)
}

func QueryTWithID[T any](r *http.Request, key string) (uuid.UUID, T, error) {
	var v T
	id, err := QueryWithID(r, key, &v)
	return id, v, err
}

func QueryWithID(r *http.Request, key string, v any) (uuid.UUID, error) {
	id, err := ID(r, key)
	if err != nil {
		return uuid.Nil, err
	}

	err = Query(r, v)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
