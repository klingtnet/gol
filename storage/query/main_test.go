package query

import (
	"errors"
	"testing"

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
	tu.ExpectEqual(t ,q.Matches[0], fields[0])
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
