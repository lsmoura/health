package dbfieldvalues

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func fieldToTags(field reflect.StructField) (string, string, bool) {
	tag := field.Tag.Get("db")
	if tag == "" {
		return "", "", false
	}
	parts := strings.Split(tag, ",")
	if len(parts) == 1 {
		return parts[0], "", true
	}
	return parts[0], parts[1], true
}

func Fields(in any, omitFields ...string) []string {
	if in == nil {
		return nil
	}

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	omitMap := make(map[string]any)
	for _, omitField := range omitFields {
		omitMap[omitField] = nil
	}

	t := v.Type()
	var fields []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name

		if field.Anonymous {
			fields = append(fields, Fields(v.Field(i).Interface(), omitFields...)...)
			continue
		}

		if name, options, ok := fieldToTags(field); ok {
			if name == "-" {
				continue
			}

			if strings.Contains(options, "inline") {
				fields = append(fields, Fields(v.Field(i).Interface(), omitFields...)...)
				continue
			}

			if name != "" {
				fieldName = name
			}
		}

		if _, ok := omitMap[fieldName]; ok {
			continue
		}
		fields = append(fields, fieldName)
	}

	return fields
}

func Values(in any, omitFields ...string) ([]any, error) {
	if in == nil {
		return nil, nil
	}

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("dbfieldvalues: expected struct, got %T", in)
	}

	omitMap := make(map[string]any)
	for _, omitField := range omitFields {
		omitMap[omitField] = nil
	}

	t := v.Type()
	var values []any
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		vField := v.Field(i)
		fieldName := field.Name

		value := vField.Interface()

		if field.Anonymous {
			innerValues, err := Values(vField.Interface())
			if err != nil {
				return nil, fmt.Errorf("error parsing %v: %w", field.Name, err)
			}
			values = append(values, innerValues...)
			continue
		}

		if name, options, ok := fieldToTags(field); ok {
			if name == "-" {
				continue
			}

			if strings.Contains(options, "inline") {
				innerValues, err := Values(vField.Interface())
				if err != nil {
					return nil, fmt.Errorf("parsing %v: %w", field.Name, err)
				}

				values = append(values, innerValues...)
				continue
			}
			if strings.Contains(options, "omitempty") && vField.IsZero() {
				value = nil
			}
			if strings.Contains(options, "json") {
				marshalledValue, err := json.Marshal(value)
				if err != nil {
					return nil, fmt.Errorf("json.Marshal %v: %w", field.Name, err)
				}

				value = marshalledValue
			}

			if name != "" {
				fieldName = name
			}
		}

		if _, ok := omitMap[fieldName]; ok {
			continue
		}

		values = append(values, value)
	}

	return values, nil
}
