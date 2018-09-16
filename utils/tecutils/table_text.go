package tecutils

import (
	"bytes"
	"fmt"
	"strings"
)

//given an array of cells, return as formatted text
//each element of body is a row and each cell has a separator
func TableToText(headers, body []string, separator string) (string, error) {

	headercount := len(headers)
	var (
		widths []int
		res    bytes.Buffer
		ast    string = "+"
	)

	for i := range headers {
		widths = append(widths, len(headers[i]))
	}

	for i := range body {
		cols := strings.Split(body[i], separator)
		if headercount != len(cols) {
			return "", fmt.Errorf("header column count don't match column count in body")
		}
		for j := range cols {
			if widths[j] < len(cols[j]) {
				widths[j] = len(cols[j])
			}
		}

	}

	print := func(elements []string, ws []int, a string, buf *bytes.Buffer, space string) {
		for i := range elements {
			buf.WriteString(fmt.Sprintf("+%s", space))
			buf.WriteString(fmt.Sprintf("%s", elements[i]) + strings.Repeat(space, ws[i]-len(elements[i])))
			buf.WriteString(space)

			if i == len(elements)-1 {
				buf.WriteString("+")
			}
		}
	}

	printline := func(buf *bytes.Buffer, ws []int, a string) {
		var e []string
		e = make([]string, len(ws))
		for i := range ws {
			e[i] = strings.Repeat(a, ws[i])
		}
		print(e, ws, a, buf, a)
		buf.WriteString("\n")
	}

	printline(&res, widths, ast)
	print(headers, widths, ast, &res, " ")
	res.WriteString("\n")
	printline(&res, widths, ast)

	for i := range body {
		row := body[i]
		print(strings.Split(row, separator), widths, ast, &res, " ")
		res.WriteString("\n")
	}

	printline(&res, widths, ast)

	return res.String(), nil
}
