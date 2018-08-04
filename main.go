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

func possibleNewVersions(ver string) []string {
	const major, minor, patch = 0, 1, 2
	v := regexp.MustCompile(`\d+`).FindAllString(ver, -1)
	switch len(v) {
	case 1:
		return []string{bump(v[major])}
	case 2:
		return []string{
			bump(v[major]) + ".0",
			v[major] + "." + bump(v[minor]),
		}
	case 3:
		return []string{
			bump(v[major]) + ".0.0",
			v[major] + "." + bump(v[minor]) + ".0",
			v[major] + "." + v[minor] + "." + bump(v[patch]),
		}
	}
	return nil
}

func bump(c string) string {
	val, err := strconv.Atoi(c)
	if err != nil {
		log.Fatal(err)
	}
	return strconv.Itoa(val + 1)
}

func latest(v string, tags map[string]struct{}) string {
	for _, p := range possibleNewVersions(v) {
		if _, ok := tags[p]; ok {
			return latest(p, tags)
		}
	}
	return v
}

func extractVersions(s string) []string {
	matches := regexp.MustCompile(`([\d\.]+)`).FindAllString(s, -1)
	return matches
}

func replaceVersions(input string) (results []string) {
	versions := extractVersions(input)
	for _, v := range versions {
		for _, nv := range possibleNewVersions(v) {
			results = append(results, strings.Replace(input, v, nv, 1))
		}
	}
	return results
}
