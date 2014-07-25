package flatjson_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/pushrax/flatjson"
)

type Child struct {
	C int    `json:"CC"`
	D string `json:"CD"`
}

func TestBasicFlatten(t *testing.T) {
	val := &struct {
		A int
		B string
	}{10, "str"}

	expected := flatjson.Map{
		"A": 10.0, // JSON numbers are all float64.
		"B": "str",
	}

	testFlattening(t, val, expected)
}

func TestEmbeddedFlatten(t *testing.T) {
	val := &struct {
		Child       // Embedded.
		Other Child // Regular child.
		A     int
	}{}

	expected := flatjson.Map{
		"A":        0.0,
		"CC":       0.0,
		"CD":       "",
		"Other.CC": 0.0,
		"Other.CD": "",
	}

	testFlattening(t, val, expected)
}

func TestIndirection(t *testing.T) {
	o2 := &Child{5, "6"}

	val := &struct {
		*Child
		Other1 interface{} `json:"O1"`
		Other2 **Child     `json:"O2"`
		Other3 *Child      `json:",omitempty"`
	}{
		Child:  &Child{1, "2"},
		Other1: &Child{3, "4"},
		Other2: &o2,
	}

	expected := flatjson.Map{
		"CC":    1.0,
		"CD":    "2",
		"O1.CC": 3.0,
		"O1.CD": "4",
		"O2.CC": 5.0,
		"O2.CD": "6",
	}

	testFlattening(t, val, expected)
}

func testFlattening(t *testing.T, val interface{}, expected flatjson.Map) {
	flat := flatjson.Flatten(val)

	enc, err := json.Marshal(flat)
	if err != nil {
		t.Fatal(err)
	}

	got := flatjson.Map{}
	err = json.Unmarshal(enc, &got)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Unmarshalled to unexpected value:\n     got: %#v\nexpected: %#v\n", got, expected)
	}
}
