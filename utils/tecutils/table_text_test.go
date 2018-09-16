package tecutils

import (
	"testing"
)

func TestTableText(t *testing.T) {
	t.Log("given a table structure, parse it into readable text")

	headers := []string{"Genre", "Decade", "Band"}
	body := []string{"Rock|80's|Rush", "Electronic|90's|Depeche Mode", "Soul|00's|Amy Winehouse", "Soul|70's|Stevie Wonder"}
	expected :=
		`
	+++++++++++++++++++++++++++++++++++++++
	+ Genre      + Decade + Band          +
	+++++++++++++++++++++++++++++++++++++++
	+ Rock       + 80's   + Rush          +
	+ Electronic + 90's   + Depeche Mode  +
	+ Soul       + 00's   + Amy Winehouse +
	+ Soul       + 70's   + Stevie Wonder +
	+++++++++++++++++++++++++++++++++++++++
	`

	result, err := TableToText(headers, body, "|")
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("\n%s", result)
	t.Log(expected)
}
