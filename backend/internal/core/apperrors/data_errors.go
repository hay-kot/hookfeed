package apperrors

import "github.com/google/uuid"

// NotFoundError is returned when a requested entity is not found. This is used in place of the
// ent.NotFoundError as we move towards SQLC and away from ent.
type NotFoundError struct {
	Property string
	Value    string
	Entity   string
}

func NewNotFoundErrorID(id uuid.UUID, entity string) NotFoundError {
	return NewNotFoundError("id", entity, id.String())
}

func NewNotFoundError(property, entity, value string) NotFoundError {
	return NotFoundError{
		Property: property,
		Value:    value,
		Entity:   entity,
	}
}

func IsNotFound(e error) bool {
	return As[NotFoundError](e)
}

func (e NotFoundError) Error() string {
	return e.Entity + " with " + e.Property + "=" + e.Value + " not found"
}
