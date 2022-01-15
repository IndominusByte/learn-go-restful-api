package validation

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

// Set a Decoder instance as a package global, because it caches
// meta-data about structs, and an instance can be shared safely.
var decoder = schema.NewDecoder()

func ParseRequest(dst interface{}, src map[string][]string) error {
	return decoder.Decode(dst, src)
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

// struct for error instance validator
type ErrorValidate struct {
	Kind, Tag, Param string
}

func StructValidate(body interface{}) map[string]interface{} {
	validate = validator.New()
	err, data := validate.Struct(body), make([]map[string]interface{}, 0)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			result := make(map[string]interface{})

			resultString := MappingMap(strings.Split(err.Namespace(), ".")[1:])
			resultSlice := strings.Split(resultString, ":")
			resultSliceLen := len(resultSlice)

			AddToMap(ErrorValidate{
				Kind:  err.Kind().String(),
				Tag:   err.Tag(),
				Param: err.Param(),
			}, result, resultSlice, 0, resultSliceLen)

			data = append(data, result)
		}

		result := make(map[string]interface{})
		for _, value := range data {
			MergeMap(value, result)
		}

		return result
	}

	return nil
}

func MergeMap(from, to map[string]interface{}) {
	for k, v := range from {
		if _, ok := to[k]; !ok {
			to[k] = v
		} else {
			MergeMap(from[k].(map[string]interface{}), to[k].(map[string]interface{}))
		}
	}
}

var insideBracket = regexp.MustCompile(`\[(.*?)\]`)
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func AddToMap(ev ErrorValidate, result map[string]interface{}, payload []string, index, length int) {
	if index+1 != length {
		p := make(map[string]interface{})
		result[payload[index]] = p
		AddToMap(ev, p, payload, index+1, length)
	} else {
		// Customize error message validator
		result[payload[index]] = SetErrorMessage(ev)
	}
}

func MappingMap(namespace []string) string {
	var mapping string
	length := len(namespace)
	for i := 0; i < length; i++ {
		bracket := insideBracket.FindString(namespace[i])
		if len(bracket) > 0 {
			namespaceCopy := insideBracket.ReplaceAllString(namespace[i], "")
			mapping += ToSnakeCase(namespaceCopy) + ":"
			mapping += ToSnakeCase(bracket[1:len(bracket)-1]) + ":"
		} else {
			mapping += ToSnakeCase(namespace[i]) + ":"
		}
	}

	return mapping[:len(mapping)-1]
}

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
