package parser

import (
	"regexp"
	"strconv"
)

func ParsePathToInt(regex, path string) (int, error) {
	re, err := regexp.Compile(regex)
	if err != nil {
		return 0, err
	}
	values := re.FindStringSubmatch(path)
	result, err := strconv.Atoi(values[1])
	if err != nil {
		return 0, err
	}
	return result, nil
}
