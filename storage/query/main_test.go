package query

import (
	"errors"
	"testing"
)

func expectEqual(t *testing.T, actual, expected interface{}) {
	if actual != expected {
		t.Errorf("%#v != %#v", actual, expected)
	}
}

func TestFind(t *testing.T) {
	b := &DefaultBuilder{}
	if _, err := b.Find("invalid field", "_").Build(); err == nil {
		t.Fail()
	}

	q, err := b.Find("id", "42").Build()
	if err != nil {
		t.Fail()
	}
	if q.Find.Name != "id" {
		t.Error("'id' must be the Name")
	}
	if q.Find.Value != "42" {
		t.Error("'42' must be the Value")
	}
}

func TestStart(t *testing.T) {
	b := &DefaultBuilder{}
	q, err := b.Start(10).Build()
	if err != nil {
		t.Fail()
	}
	if q.Start != 10 {
		t.Error(q.Start, "!=", 10)
	}
}

func TestMatchSingle(t *testing.T) {
	b := &DefaultBuilder{}
	q, err := b.Match("title", "cool").Build()
	if err != nil {
		t.Fail()
	}
	if len(q.Matches) != 1 {
		t.Fatal(len(q.Matches), "!=", 1)
	}
	field := Field{"title", "cool"} 
	if q.Matches[0] != field {
		t.Error(q.Matches, "!=", field)
	}
}

func TestMatchMultiple(t *testing.T) {
	b := &DefaultBuilder{}
	q, err := b.Match("id", "42").
		Match("title", "cool").Build()
	if err != nil {
		t.Fatal("invalid .Match query:", err)
	}

	if len(q.Matches) != 2 {
		t.Fatal(len(q.Matches), "!=", 2)
	}
}

func TestValueInHelper(t *testing.T) {
	allowed := []string{"oops"}
	if err := valueIn("_", "hey", allowed); err == nil {
		t.Error("hey is not an allowed value")
	}

	if err := valueIn("_", "oops", allowed); err != nil {
		t.Error("oops is an allowed value")
	}
}

func TestInvalidBuildMustError(t *testing.T) {
	b := Invalid{errors.New("oops")}
	if q, err := b.Build(); q != nil || err == nil {
		t.Fail()
	}
}
