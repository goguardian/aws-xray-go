package utils

import (
	"os"
	"testing"
)

func TestGetenvOrDefault(t *testing.T) {
	key := "NOT_DEFINED"
	def := "default"

	if val := GetenvOrDefault(key, def); val != def {
		t.Error("Expected default value")
	}

	os.Setenv(key, "not default")

	if val := GetenvOrDefault(key, def); val == def {
		t.Error("Expected non-default value")
	}
}
