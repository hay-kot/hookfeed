package apperrors_test

import (
	"testing"

	"github.com/hay-kot/hookfeed/backend/internal/core/apperrors"
)

type errPointer struct{}

func (*errPointer) Error() string {
	return "errPointer"
}

func Test_As_PointerType(t *testing.T) {
	t.Parallel()

	var err error = &errPointer{}

	if !apperrors.As[*errPointer](err) {
		t.Error("expected true")
	}
}

type errValue struct{}

func (errValue) Error() string {
	return "errPointer"
}

func Test_As_ValueType(t *testing.T) {
	t.Parallel()

	var err error = errValue{}

	if !apperrors.As[errValue](err) {
		t.Error("expected true")
	}
}
