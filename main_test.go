package main

import (
	"testing"
)

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

func TestLatest(t *testing.T) {
	tags := map[string]struct{}{
		"latest": struct{}{},
		"1":      struct{}{},
		"1.10":   struct{}{},
		"1.11.0": struct{}{},
	}
	input := "1.10.2"
	want := "1.11.0"
	got := latest(input, tags)
	if got != want {
		t.Errorf("latest(%q) = %q, want %q", input, got, want)
	}
}

func TestExtractVersions(t *testing.T) {
	tcs := []struct {
		input string
		wants []string
	}{
		{
			input: "1.10.3-alpine3.7",
			wants: []string{
				"1.10.3",
				"3.7",
			},
		},
	}
	for _, tc := range tcs {
		got := extractVersions(tc.input)
		for i, want := range tc.wants {
			if got[i] != want {
				t.Errorf("extractVersions(%q) = %q, want %q", tc.input, got[i], want)
			}
		}
	}

}

func TestReplaceVersions(t *testing.T) {
	input := "1.10.3-foo2.1"
	wants := []string{
		"2.0.0-foo2.1",
		"1.11.0-foo2.1",
		"1.10.4-foo2.1",
		"1.10.3-foo3.0",
		"1.10.3-foo2.2",
	}
	got := replaceVersions(input)
	for i, want := range wants {
		if want != got[i] {
			t.Errorf("replaceVersions() = %q, want %q", got, want[i])
		}
	}
}

func TestParseTags(t *testing.T) {
	data := []byte(`[{"layer": "", "name": "latest"}, {"layer": "", "name": "1"}, {"layer": "", "name": "1-alpine"}, {"layer": "", "name": "1-alpine-perl"}, {"layer": "", "name": "1-perl"}, {"layer": "", "name": "1.10"}, {"layer": "", "name": "1.10-alpine"}]`)
	wants := []string{
		"latest",
		"1",
		"1-alpine",
		"1-alpine-perl",
		"1-perl",
		"1.10",
		"1.10-alpine",
	}
	got := parseTags(data)
	for i, want := range wants {
		if got[i] != want {
			t.Errorf("parseTags() = %q, want %q", got[i], want)
		}
	}
}
