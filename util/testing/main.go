// some utilities to make tests easier
package testing

import (
	"reflect"
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

func ExpectNil(t *testing.T, value interface{}) {
	if !isNil(value) {
		t.Errorf("%#v should be nil", value)
	}
}

func RequireNil(t *testing.T, value interface{}) {
	if !isNil(value) {
		t.Fatalf("%#v should be nil", value)
	}
}

func ExpectNotNil(t *testing.T, value interface{}) {
	if isNil(value) {
		t.Errorf("%#v should be non-nil", value)
	}
}

func RequireNotNil(t *testing.T, value interface{}) {
	if isNil(value) {
		t.Fatalf("%#v should be non-nil", value)
	}
}

func isNil(value interface{}) bool {
	return value == nil || reflect.ValueOf(value).IsNil()
}
