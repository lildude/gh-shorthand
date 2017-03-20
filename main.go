package main

import (
	"encoding/json"
	"fmt"
	"github.com/zerowidth/alfred-gh-shorthand/alfred"
	"os"
)

func main() {
	var input string
	if len(os.Args) < 2 {
		input = ""
	} else {
		input = os.Args[1]
	}
	fmt.Fprintf(os.Stderr, "input: %#v\n", input)
	item := alfred.Item{
		Title: "hello",
		Valid: false,
	}
	items := alfred.Items{Items: []alfred.Item{item}}
	encoded, _ := json.Marshal(items)
	os.Stdout.Write(encoded)
	os.Stdout.WriteString("\n")
}
