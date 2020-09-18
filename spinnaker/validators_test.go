package spinnaker

import (
	"testing"
)

func TestValidateApplicationName(t *testing.T) {
	validNames := []string{
		"ValidName",
		"validname",
		"invalid-name",
	}
	for _, v := range validNames {
		_, errors := validateApplicationName(v, "application")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Application name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"invalid:name",
		"invalid name",
		"invalid_name",
		"",
	}
	for _, v := range invalidNames {
		_, errors := validateApplicationName(v, "application")
		if len(errors) == 0 {
			t.Fatalf("%q should be a valid Application name", v)
		}
	}
}
