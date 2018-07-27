package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(replace(string(data)))
}

func replace(data string) string {
	re := regexp.MustCompile(`(from|FROM) (\S+):(\S+)`)
	return re.ReplaceAllString(data, "$1 $2:latest")
}
