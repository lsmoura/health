package dbfieldvalues

import (
	"reflect"
	"testing"
)

func TestFields(t *testing.T) {
	type FooStruct struct {
		Foo string `db:"foo"`
	}
	tests := []struct {
		in     any
		fields []string
		omit   []string
	}{
		{
			in:     nil,
			fields: nil,
		},
		{
			in: struct {
				A string
			}{},
			fields: []string{"A"},
		},
		{
			in: struct {
				A string `db:"a"`
			}{},
			fields: []string{"a"},
		},
		{
			in: struct {
				A string `db:"-"`
			}{},
			fields: nil,
		},
		{
			// anonymous field
			in: struct {
				FooStruct
				A string `db:"a"`
			}{},
			fields: []string{"foo", "a"},
		},
		{
			// inline field
			in: struct {
				Foo FooStruct `db:",inline"`
				B   string    `db:"b"`
			}{},
			fields: []string{"foo", "b"},
		},
		{
			// omit field
			in: struct {
				A string `db:"a"`
				B string `db:"b"`
			}{},
			omit:   []string{"b"},
			fields: []string{"a"},
		},
	}

	for _, test := range tests {
		fields := Fields(test.in, test.omit...)
		if !reflect.DeepEqual(fields, test.fields) {
			t.Errorf("Fields(%#v): expected %#v, got %#v", test.in, test.fields, fields)
		}
	}
}

func TestValues(t *testing.T) {
	type FooStruct struct {
		Foo string `db:"foo" json:"foo"`
	}

	tests := []struct {
		in     any
		omit   []string
		values []any
	}{
		{
			in:     nil,
			values: nil,
		},
		{
			in: struct {
				A string
			}{A: "aa"},
			values: []any{"aa"},
		},
		{
			in: struct {
				A string `db:"a"`
				B string `db:"b"`
				C string `db:"c"`
			}{A: "aa", B: "bb", C: "cc"},
			omit:   []string{"b"},
			values: []any{"aa", "cc"},
		},
		{
			// anonymous field
			in: struct {
				FooStruct
				A string `db:"a"`
			}{
				FooStruct: FooStruct{Foo: "fooValue"},
				A:         "aValue",
			},
			values: []any{"fooValue", "aValue"},
		},
		{
			// inline field
			in: struct {
				Foo FooStruct `db:",inline"`
				B   string    `db:"b"`
			}{
				Foo: FooStruct{Foo: "fooValue"},
				B:   "bValue",
			},
			values: []any{"fooValue", "bValue"},
		},
		{
			// empty field
			in: struct {
				B string `db:"b"`
			}{
				B: "",
			},
			values: []any{""},
		},
		{
			// omitempty field
			in: struct {
				B string `db:"b,omitempty"`
			}{
				B: "",
			},
			values: []any{nil},
		},
		{
			// json field
			in: struct {
				Foo FooStruct `db:"foo,json"`
			}{
				Foo: FooStruct{Foo: "fooValue"},
			},
			values: []any{[]byte("{\"foo\":\"fooValue\"}")},
		},
	}

	for _, test := range tests {
		values, err := Values(test.in, test.omit...)
		if err != nil {
			t.Errorf("Values(%#v): expected %#v, got %#v", test.in, test.values, values)
			continue
		}
		if !reflect.DeepEqual(values, test.values) {
			t.Errorf("Values(%#v): expected %#v, got %#v", test.in, test.values, values)
		}
	}
}
