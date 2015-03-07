// some utilities to make tests easier
package testing

import (
	"testing"
)

func ExpectEqual(t *testing.T, actual, expected interface{}) {
	if actual != expected {
		t.Errorf("%#v != %#v", actual, expected)
	}
}

func RequireEqual(t *testing.T, actual, expected interface{}) {
	if actual != expected {
		t.Fatalf("%#v != %#v", actual, expected)
	}
}

func ExpectNotNil(t *testing.T, value interface{}) {
	if value == nil {
		t.Error("%#v should be non-nil")
	}
}

func RequireNotNil(t *testing.T, value interface{}) {
	if value == nil {
		t.Fatal("%#v should be non-nil")
	}
}
