package main

import "testing"

func TestReplace(t *testing.T) {
	input := `
from golang:1.10.2-alpine3.7 AS build
WORKDIR /
FROM foobar
FROM scratch:1.10 as foo
`
	want := `
from golang:latest AS build
WORKDIR /
FROM foobar
FROM scratch:latest as foo
`
	got := replace(input)
	if got != want {
		t.Errorf("replace() =\n%s\nwant:\n%s", got, want)
	}
}
