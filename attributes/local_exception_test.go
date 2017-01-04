package attributes

import (
	"errors"
	"testing"
)

func TestNewLocalException(t *testing.T) {
	err := errors.New("An error")
	exception := NewLocalException(err)

	if exception.Message != err.Error() {
		t.Error("Local exception exception message incorrect")
	}
}
