package query

import (
	"errors"
	"net/url"
	"testing"
	"time"

	tu "../../util/testing"
)

func TestIsDefault(t *testing.T) {
	tu.ExpectEqual(t, IsDefault(Default), true)
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
	tu.ExpectEqual(t, q.Find.Name, "id")
	tu.ExpectEqual(t, q.Find.Value, "42")
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

	tu.RequireEqual(t, len(q.Matches), 1)

	field := Field{"title", "cool"}
	tu.ExpectEqual(t, q.Matches[0], field)
}

func TestMatchMultiple(t *testing.T) {
	b := &DefaultBuilder{}
	q, err := b.Match("id", "42").
		Match("title", "cool").Build()
	if err != nil {
		t.Fatal("invalid .Match query:", err)
	}

	tu.RequireEqual(t, len(q.Matches), 2)

	fields := []Field{Field{"id", "42"}, Field{"title", "cool"}}
	tu.ExpectEqual(t, q.Matches[0], fields[0])
	tu.ExpectEqual(t, q.Matches[1], fields[1])
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

func TestFromParamsFind(t *testing.T) {
	q, _ := fromParams(t, "http://not.es/find?id=1")

	tu.RequireNotNil(t, q.Find)
	tu.ExpectEqual(t, *q.Find, Field{"id", "1"})
}

func TestFromParamsFindMultiple(t *testing.T) {
	q, _ := fromParams(t, "http://not.es/find?id=1&title=Hello, World!")

	tu.RequireNotNil(t, q.Find)
	tu.ExpectEqual(t, *q.Find, Field{"title", "Hello, World!"})
}

func TestFromParamsFindMultipleSame(t *testing.T) {
	q, _ := fromParams(t, "http://not.es/find?id=1&id=2")
	tu.RequireNotNil(t, q.Find)
	tu.ExpectEqual(t, *q.Find, Field{"id", "2"})

	q, _ = fromParams(t, "http://not.es/find?title=hey&title=ho")
	tu.RequireNotNil(t, q.Find)
	tu.ExpectEqual(t, *q.Find, Field{"title", "ho"})
}

func TestFromParamsStart(t *testing.T) {
	q, _ := fromParams(t, "http://not.es/find?start=42")
	tu.ExpectEqual(t, q.Start, 42)
}

func TestFromParamsCount(t *testing.T) {
	q, _ := fromParams(t, "http://not.es/find?count=42")
	tu.ExpectEqual(t, q.Count, 42)
}

func TestFromParamsSortBy(t *testing.T) {
	q, _ := fromParams(t, "http://not.es/find")
	tu.ExpectEqual(t, q.SortBy, "created")

	q, _ = fromParams(t, "http://not.es/find?sort=created")
	tu.ExpectEqual(t, q.SortBy, "created")

	q, _ = fromParams(t, "http://not.es/find?sort=title")
	tu.ExpectEqual(t, q.SortBy, "title")

	// invalid params checked by fromParams
}

func TestFromParamsReverse(t *testing.T) {
	q, _ := fromParams(t, "http://not.es/find?reverse")
	tu.ExpectEqual(t, q.Reverse, true)

	q, _ = fromParams(t, "http://not.es/find?reverse=true")
	tu.ExpectEqual(t, q.Reverse, true)

	q, _ = fromParams(t, "http://not.es/find?reverse=false")
	tu.ExpectEqual(t, q.Reverse, false)
}

func TestFromParamsMatch(t *testing.T) {
	q, _ := fromParams(t, "http://not.es/find")
	tu.RequireEqual(t, len(q.Matches), 0)

	q, _ = fromParams(t, "http://not.es/find?match=id:x")
	tu.RequireEqual(t, len(q.Matches), 1)
	tu.ExpectEqual(t, q.Matches[0], Field{"id", "x"})

	q, _ = fromParams(t, "http://not.es/find?match=title:cool")
	tu.RequireEqual(t, len(q.Matches), 1)
	tu.ExpectEqual(t, q.Matches[0], Field{"title", "cool"})

	q, _ = fromParams(t, "http://not.es/find?match=title:cool&match=content:very")
	tu.RequireEqual(t, len(q.Matches), 2)
	tu.ExpectEqual(t, q.Matches[0], Field{"title", "cool"})
	tu.ExpectEqual(t, q.Matches[1], Field{"content", "very"})
}

func TestFromParamsRange(t *testing.T) {
	q, _ := fromParams(t, "http://not.es/find?range=2015-03-03T09:35:00Z,2015-03-09T10:50:00Z")
	tu.RequireNotNil(t, q.RangeStart)
	tu.ExpectEqual(t, *q.RangeStart, time.Date(2015, 3, 3, 9, 35, 0, 0, time.UTC))

	tu.RequireNotNil(t, q.RangeEnd)
	tu.ExpectEqual(t, *q.RangeEnd, time.Date(2015, 3, 9, 10, 50, 0, 0, time.UTC))
}

func fromParams(t *testing.T, rawUrl string) (*Query, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	vals := u.Query()
	q, err := FromParams(vals)
	if err != nil {
		t.Fatalf("%s should be parsable: %s", rawUrl, err)
		return nil, err
	}
	return q, nil
}
