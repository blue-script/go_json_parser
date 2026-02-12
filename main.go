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

var ErrorIncorrectCommand = "Error: incorrect command"
var ErrorNoJsonData = "Error: no json data"
var ErrorInvalidJson = "Error: invalid JSON"
var ErrorKeyNotFound = "Error: key not found"
var ErrorInvalidIndex = "Error: invalid array index"
var ErrorIndexOutOfRange = "Error: index is out of range"
var ErrorNotObject = "Error: path does not lead to an object"

type JSONValue = any

type JSONParser struct {
	input []byte
}

func NewJSONParser(input []byte) *JSONParser {
	return &JSONParser{
		input: input,
	}
}

func (j *JSONParser) Parse() (JSONValue, error) {
	input := j.input
	pos := 0
	var result JSONValue
	var err error

	pos = skipSpaces(input, pos)

	if pos >= len(j.input) {
		return nil, errors.New(ErrorInvalidJson)
	}

	switch input[pos] {
	case '{':
		result, pos, err = parseObject(input, pos)
		if err != nil {
			return nil, errors.New(ErrorInvalidJson)
		}

	case '[':
		result, pos, err = parseArray(input, pos)
		if err != nil {
			return nil, errors.New(ErrorInvalidJson)
		}

	default:
		return nil, errors.New(ErrorInvalidJson)
	}

	pos = skipSpaces(input, pos)

	if pos != len(input) {
		return nil, errors.New(ErrorInvalidJson)
	}

	// temporary
	if false {
		return nil, errors.New(ErrorInvalidJson)
	}
	// temporary
	return result, nil
}

// parse {"a": {}, "b": {}}

func parseObject(data []byte, pos int) (JSONValue, int, error) {
	if pos >= len(data) {
		return nil, pos, errors.New(ErrorInvalidJson)
	}

	currentValue := map[string]any{}

	if data[pos] != '{' {
		return nil, pos, errors.New(ErrorInvalidJson)
	}

	pos++
	pos = skipSpaces(data, pos)
	if pos >= len(data) {
		return nil, pos, errors.New(ErrorInvalidJson)
	}

	if data[pos] == '}' {
		return currentValue, pos + 1, nil
	}

	for {
		if pos >= len(data) {
			return nil, pos, errors.New(ErrorInvalidJson)
		}

		if data[pos] != '"' {
			return nil, pos, errors.New(ErrorInvalidJson)
		}
		var key string
		var err error
		key, pos, err = parseString(data, pos)
		if err != nil {
			return nil, pos, err
		}

		pos = skipSpaces(data, pos)
		if pos >= len(data) {
			return nil, pos, errors.New(ErrorInvalidJson)
		}
		if data[pos] != ':' {
			return nil, pos, errors.New(ErrorInvalidJson)
		}

		pos++
		pos = skipSpaces(data, pos)
		var value JSONValue
		value, pos, err = parseValue(data, pos)
		if err != nil {
			return nil, pos, err
		}

		currentValue[key] = value

		pos = skipSpaces(data, pos)
		if pos >= len(data) {
			return nil, pos, errors.New(ErrorInvalidJson)
		}

		if data[pos] == ',' {
			pos++
			pos = skipSpaces(data, pos)
			continue
		}

		if data[pos] == '}' {
			return currentValue, pos + 1, nil
		}

		return nil, pos, errors.New(ErrorInvalidJson)
	}
}

func parseValue(data []byte, pos int) (JSONValue, int, error) {
	pos = skipSpaces(data, pos)
	if pos >= len(data) {
		return nil, pos, errors.New(ErrorInvalidJson)
	}

	var result JSONValue
	var err error

	switch data[pos] {
	case '{':
		result, pos, err = parseObject(data, pos)
		if err != nil {
			return nil, pos, errors.New(ErrorInvalidJson)
		}

	case '[':
		result, pos, err = parseArray(data, pos)
		if err != nil {
			return nil, pos, errors.New(ErrorInvalidJson)
		}

	case '"':
		result, pos, err = parseString(data, pos)
		if err != nil {
			return nil, pos, errors.New(ErrorInvalidJson)
		}

	default:
		return nil, pos, errors.New(ErrorInvalidJson)
	}

	return result, pos, nil
}

func parseString(data []byte, pos int) (string, int, error) {
	first := pos + 1
	if data[pos] != '"' {
		return "", pos, errors.New(ErrorInvalidJson)
	}
	pos++
	for ; pos < len(data); pos++ {
		if data[pos] == '"' {
			return string(data[first:pos]), pos + 1, nil
		}
	}

	return "", pos, errors.New(ErrorInvalidJson)
}

func parseArray(data []byte, pos int) (JSONValue, int, error) {
	// var currentValue JSONValue
	if pos >= len(data) {
		return nil, pos, errors.New(ErrorInvalidJson)
	}

	if data[pos] != '[' {
		return nil, pos, errors.New(ErrorInvalidJson)
	}

	pos++
	pos = skipSpaces(data, pos)
	if pos >= len(data) {
		return nil, pos, errors.New(ErrorInvalidJson)
	}

	if data[pos] == ']' {
		return []any{}, pos + 1, nil
	}

	return nil, pos, errors.New(ErrorInvalidJson)
}

func skipSpaces(data []byte, pos int) int {
Loop:
	for ; pos < len(data); pos++ {
		switch data[pos] {
		case ' ', '\n', '\t', '\r':
			continue
		default:
			break Loop
		}
	}

	return pos
}

func iterHelper(root JSONValue) (iter.Seq2[string, JSONValue], error) {
	switch v := root.(type) {
	case map[string]any:
		{
			return func(yield func(string, JSONValue) bool) {
				for key, value := range v {
					if !yield(key, value) {
						return
					}
				}
			}, nil
		}
	case []any:
		{
			return func(yield func(string, JSONValue) bool) {
				for i, value := range v {
					if !yield(strconv.Itoa(i), value) {
						return
					}
				}
			}, nil
		}
	default:
		return nil, fmt.Errorf("Error: cannot iterate type %T", v)
	}
}

func PathIterator(root any, path ...string) (iter.Seq2[string, JSONValue], error) {
	if root == nil {
		return nil, errors.New(ErrorNoJsonData)
	}

	cur := root

	for i := range len(path) {
		switch v := cur.(type) {
		case []any:
			{
				num, err := strconv.Atoi(path[i])
				if err != nil {
					return nil, fmt.Errorf("Error: invalid array index '%s'", path[i])
				}
				if num >= len(v) || num < 0 {
					return nil, fmt.Errorf("Error: index '%d' is out of range", num)
				}
				cur = v[num]
			}
		case map[string]any:
			{
				m, ok := v[path[i]]
				if !ok {
					return nil, fmt.Errorf("Error: key '%s' not found", path[i])
				}
				cur = m
			}

		default:
			return nil, errors.New(ErrorNotObject)
		}
	}

	switch cur.(type) {
	case map[string]any, []any:
		return iterHelper(cur)
	default:
		return nil, errors.New(ErrorNotObject)
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

		switch cmd {
		case "parse":
			if len(parts) < 2 {
				fmt.Println(ErrorIncorrectCommand)
				continue
			}

			arg = parts[1]

			parseJsonData, err = NewJSONParser([]byte(arg)).Parse()
			if err != nil {
				fmt.Println(ErrorInvalidJson)
				continue
			}

		case "iterate":
			if len(parts) > 1 {
				arg = parts[1]
			}

			var path []string
			if strings.TrimSpace(arg) != "" {
				path = strings.Split(arg, ".")
			}

			iterator, err := PathIterator(parseJsonData, path...)
			if err != nil {
				fmt.Println(err)
				continue
			}

			printKV := func(key string, value JSONValue) bool {
				fmt.Printf("%s: %v\n", key, value)
				return true
			}
			iterator(printKV)

		case "exit":
			break CommandLoop
		default:
			fmt.Println(ErrorIncorrectCommand)
		}
	}
}

// func test() {
// 	input := `{"foo": {"bar": {"baz": 42, "qux": "hello"}}}`
// 	// input := `[{"foo":"1"},{"foo":"2"},{"foo":"3"}]`
// 	result, err := NewJSONParser([]byte(input)).Parse()
// 	if err != nil {
// 		log.Fatalln("Error: invalid JSON")
// 	}
// 	path := []string{"foo", "bar"}
// 	// path := []string{"1"}
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
