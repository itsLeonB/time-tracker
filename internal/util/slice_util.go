package util

import (
	"reflect"

	"github.com/rotisserie/eris"
)

func MapSlice[T any, U any](input []T, mapperFunc func(T) U) []U {
	output := make([]U, len(input))

	for i, v := range input {
		output[i] = mapperFunc(v)
	}

	return output
}

func MapSliceWithError[T any, U any](input []T, mapperFunc func(T) (U, error)) ([]U, error) {
	output := make([]U, len(input))

	for i, v := range input {
		mapped, err := mapperFunc(v)
		if err != nil {
			return nil, err
		}

		output[i] = mapped
	}

	return output, nil
}

func GroupBy[T any](items []T, fieldName string) (map[any][]T, error) {
	result := make(map[any][]T)
	for _, item := range items {
		val := reflect.ValueOf(item)

		// Support pointer or value
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		if val.Kind() != reflect.Struct {
			return nil, eris.Errorf("GroupBy: expected struct but got %s", val.Kind())
		}

		field := val.FieldByName(fieldName)
		if !field.IsValid() {
			return nil, eris.Errorf("GroupBy: field %s not found", fieldName)
		}

		key := field.Interface()
		if result[key] == nil {
			result[key] = make([]T, 0)
		}

		result[key] = append(result[key], item)
	}

	return result, nil
}
