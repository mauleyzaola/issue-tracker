package tecutils

import (
	"strings"
)

//splits a string based on sep char
//removes empty elements
func Split(text string, sep string) []string {
	values := strings.Split(text, sep)
	var res []string
	for _, v := range values {
		if len(v) != 0 {
			res = append(res, v)
		}
	}
	return res
}

func SplitUnique(text string, sep string) []string {
	values := Split(text, sep)
	m := make(map[string]string)
	for _, v := range values {
		m[v] = v
	}
	values = make([]string, 0)
	for i := range m {
		values = append(values, i)
	}
	return values
}
