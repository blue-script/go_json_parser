package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type JSONParser struct {
	input string
}

func (j *JSONParser) Parse() (any, error) {

	return map[string]any{"foo": "hello"}, nil
}

func NewJSONParser(input string) *JSONParser {
	return &JSONParser{
		input: input,
	}
}

func iterate() {}

func main() {
	// test()

	var parseJsonData any
	var err error

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		// cmd := strings.ToLower(parts[0])
		// arg := ""
		if len(parts) < 2 {
			fmt.Println("Incorrect command")
			continue
		}
		switch strings.ToLower(parts[0]) {
		case "parse":
			parseJsonData, err = NewJSONParser(parts[1]).Parse()
			fmt.Println(parseJsonData)
			if err != nil {
				fmt.Println("Error parsing JSON")
				return
			}
		case "iterate":
		default:
			fmt.Println("Incorrect command")
		}
	}
}

// func test() {
// 	input := `{"foo": {"bar": {"baz": 42, "qux": "hello"}}}`
// 	result, err := NewJSONParser(input).Parse()
// 	if err != nil {
// 		log.Fatalln("Error: invalid JSON")
// 	}
// 	path := []string{"foo", "bar"}
// 	iter, err := PathIterator(result, path...)
// 	if err != nil {
// 		log.Fatalln("Error:", err)
// 	}
// 	printKV := func(key string, value JSONValue) bool {
// 		log.Printf("%s: %v\n", key, value)
// 		return true
// 	}
// 	iter(printKV)
// }
