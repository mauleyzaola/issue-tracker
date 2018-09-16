package tecutils

import (
	"testing"
	"time"
)

func TestCloneStruct(t *testing.T) {
	type nested struct {
		NestedProp string
	}
	type myStruct struct {
		IField   int
		SField   string
		DField   time.Time
		BolField bool
		I2Field  int64
		FField   float64
		Nested   *nested
	}

	source := &myStruct{}
	source.BolField = true
	source.DField = time.Now()
	source.FField = 1811.55
	source.I2Field = 45646546546546556
	source.IField = 255
	source.Nested = &nested{NestedProp: "some value"}
	source.SField = "another1"

	_, err := Clone(source, nil)
	if err == nil {
		t.Error()
	}

	result, err := Clone(source, &myStruct{})
	if err != nil {
		t.Error(err)
	}
	target, ok := result.(*myStruct)
	if !ok {
		t.Error()
	}
	if target == nil {
		t.Error()
	}

	if source.BolField != target.BolField {
		t.Error()
	}
	if source.DField.Second() != target.DField.Second() {
		t.Error()
	}
	if source.FField != target.FField {
		t.Error()
	}
	if source.I2Field != target.I2Field {
		t.Error()
	}
	if source.IField != target.IField {
		t.Error()
	}
	if source.SField != target.SField {
		t.Error()
	}
}
