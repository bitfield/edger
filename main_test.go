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
	got := string(replace([]byte(input)))
	if got != want {
		t.Errorf("replace() =\n%s\nwant:\n%s", got, want)
	}
}

func TestPossibleNewVersions(t *testing.T) {
	tcs := []struct {
		input string
		wants []string
	}{
		{
			input: "6",
			wants: []string{
				"7",
			},
		},

		{
			input: "3.7",
			wants: []string{
				"4.0",
				"3.8",
			},
		},
		{
			input: "3.7.4",
			wants: []string{
				"4.0.0",
				"3.8.0",
				"3.7.5",
			},
		},
	}
	for _, tc := range tcs {
		got := possibleNewVersions(tc.input)
		for i, want := range tc.wants {
			if got[i] != want {
				t.Errorf("possibleNewVersions(%q) = %q, want %q", tc.input, got[i], want)
			}
		}
	}
}
