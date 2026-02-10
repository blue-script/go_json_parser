package main

import (
	"bufio"
	"errors"
	"fmt"
	"iter"
	"os"
	"strconv"
	"strings"
)

type JSONParser struct {
	input []byte
}

func (j *JSONParser) Parse() (any, error) {

	return map[string]any{"foo": "hello"}, nil
}

func NewJSONParser(input []byte) *JSONParser {
	return &JSONParser{
		input: input,
	}
}

func Iterate(root any, path ...string) (iter.Seq2[any, any], error) {
	if root == nil {
		return nil, errors.New("No json data")
	}

	if len(path) == 0 {
		switch v := root.(type) {
		case map[string]any:
			{
				return func(yield func(any, any) bool) {
					for key, value := range v {
						if !yield(key, value) {
							return
						}
					}
				}, nil
			}
		case []any:
			{
				return func(yield func(any, any) bool) {
					for i, value := range v {
						if !yield(i, value) {
							return
						}
					}
				}, nil
			}
		default:
			return nil, fmt.Errorf("Cannot iterate type %T\n", v)
		}
	}

	cur := root
	for i := range len(path) {
		switch v := cur.(type) {
		case []any:
			{
				num, err := strconv.Atoi(path[i])
				if err != nil {
					return nil, fmt.Errorf("invalid array index %q", path[i])
				}
				cur = v[num]
			}
		case map[string]any:
			{
				m, ok := v[path[i]]
				if !ok {
					return nil, fmt.Errorf("There isn't this key %s in map", path[i])
				}
				cur = m
			}

		default:
			err := fmt.Sprintf("Cannot iterate type %T\n", cur)
			return nil, errors.New(err)
		}
	}

	switch v := cur.(type) {
	case map[string]any:
		{
			return func(yield func(any, any) bool) {
				for key, value := range v {
					if !yield(key, value) {
						return
					}
				}
			}, nil
		}
	case []any:
		{
			return func(yield func(any, any) bool) {
				for i, value := range v {
					if !yield(i, value) {
						return
					}
				}
			}, nil
		}
	default:
		return nil, fmt.Errorf("Cannot iterate type %T\n", v)
	}
}

func main() {
	// test()

	var parseJsonData any
	var err error

	scanner := bufio.NewScanner(os.Stdin)
CommandLoop:
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		cmd := strings.ToLower(parts[0])
		arg := ""

		if cmd != "exit" && len(parts) < 2 {
			fmt.Println("Incorrect command")
			continue
		}

		switch cmd {
		case "parse":
			parseJsonData, err = NewJSONParser([]byte(arg)).Parse()
			fmt.Println(parseJsonData)
			if err != nil {
				fmt.Println("Error parsing JSON")
				return
			}
		case "iterate":
			Iterate(
				map[string]any{
					"foo": map[string]any{
						"bar": map[string]any{
							"baz": 42,
							"qux": "hello",
						},
					},
				},
				"foo", "bar",
			)
		case "exit":
			break CommandLoop
		default:
			fmt.Println("Incorrect command")
		}
	}
}

// func test() {
// 	input := `{"foo": {"bar": {"baz": 42, "qux": "hello"}}}`
// 	input := `[{"foo":"1"},{"foo":"2"},{"foo":"3"}]`
// 	result, err := NewJSONParser(input).Parse()
// 	if err != nil {
// 		log.Fatalln("Error: invalid JSON")
// 	}
// 	path := []string{"foo", "bar"}
// 	path := []string{"1"}
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
