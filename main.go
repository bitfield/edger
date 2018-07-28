package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	dockerfile := flag.String("f", "", "Source Dockerfile")
	flag.Parse()
	if *dockerfile == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		*dockerfile = filepath.Join(wd, "Dockerfile")
	}
	data, err := ioutil.ReadFile(*dockerfile)
	if err != nil {
		log.Fatal(err)
	}
	outFile := *dockerfile + ".edge"
	if err = ioutil.WriteFile(outFile, replace(data), 0644); err != nil {
		log.Fatal(err)
	}
}

func replace(data []byte) []byte {
	re := regexp.MustCompile(`(from|FROM) (\S+):(\S+)`)
	return re.ReplaceAll(data, []byte("$1 $2:latest"))
}

type version struct {
	major string
	minor string
	patch string
}

func (v version) String() string {
	var parts []string
	if v.major != "" {
		parts = append(parts, v.major)
	}
	if v.minor != "" {
		parts = append(parts, v.minor)

	}
	if v.patch != "" {
		parts = append(parts, v.patch)
	}
	return strings.Join(parts, ".")
}

func possibleNewVersions(ver string) (versions []string) {
	var major, minor, patch string
	re := regexp.MustCompile(`(\d+)`)
	components := re.FindAllString(ver, -1)
	switch len(components) {
	case 1:
		major = components[0]
		versions = append(versions, bump(major))
	case 2:
		major = components[0]
		versions = append(versions, version{
			major: bump(major),
			minor: "0",
		}.String())
		minor = components[1]
		versions = append(versions, version{
			major: major,
			minor: bump(minor),
		}.String())
	case 3:
		major = components[0]
		minor = components[1]
		patch = components[2]
		versions = append(versions, version{
			major: bump(major),
			minor: "0",
			patch: "0",
		}.String())
		versions = append(versions, version{
			major: major,
			minor: bump(minor),
			patch: "0",
		}.String())
		versions = append(versions, version{
			major: major,
			minor: minor,
			patch: bump(patch),
		}.String())
	}
	return versions
}

func bump(c string) string {
	val, err := strconv.Atoi(c)
	if err != nil {
		log.Fatal(err)
	}
	return strconv.Itoa(val + 1)
}
