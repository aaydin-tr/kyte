package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var builtinTags = []string{"omitempty", "minsize", "truncate", "inline"}

var (
	ErrNilSource    = errors.New("source is nil")
	ErrNotPtrSource = errors.New("source is not a pointer")

	ErrEmptyField             = errors.New("field is empty")
	ErrNilPointerField        = errors.New("field is nil pointer")
	ErrFieldMustBePtrOrString = errors.New("field must be string or pointer of an source struct field")

	ErrNotValidFieldForQuery = errors.New("this field is not in the source struct")

	ErrNilValue = errors.New("value is nil")
)

type Kyte struct {
	source     any
	fields     map[any]string
	fieldNames []string
	errs       []error
	fieldCheck bool
}

func newKyte(source any) *Kyte {
	kyte := &Kyte{fields: make(map[any]string), fieldCheck: true}
	if source == nil {
		kyte.errs = append(kyte.errs, ErrNilSource)
		return kyte
	}

	k := reflect.ValueOf(source)
	if k.Kind() != reflect.Ptr {
		kyte.errs = append(kyte.errs, ErrNotPtrSource)
		return kyte
	}

	kyte.source = source
	v := reflect.ValueOf(source).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldValue := v.Field(i)
		if field.Anonymous {
			continue
		}

		bsonTag := getBsonTag(field)
		if bsonTag != "" {
			kyte.fields[fieldValue.Addr().Interface()] = bsonTag
		}

		if field.Type.Kind() == reflect.Struct {
			getSubStructFields(v.Field(i), bsonTag+".", kyte.fields)
		}
	}

	for _, v := range kyte.fields {
		kyte.fieldNames = append(kyte.fieldNames, v)
	}

	return kyte
}

func getBsonTag(field reflect.StructField) string {
	bsonTag := field.Tag.Get("bson")
	if bsonTag == "" {
		return ""
	}

	splitBsonTag := strings.Split(bsonTag, ",")

	for _, tag := range splitBsonTag {
		if !contains(builtinTags, tag) {
			return tag
		}
	}

	return ""
}

func getSubStructFields(s reflect.Value, parentPrefix string, fields map[any]string) {
	for i := 0; i < s.NumField(); i++ {
		field := s.Type().Field(i)
		fieldValue := s.Field(i)
		if field.Anonymous {
			continue
		}

		if field.Type.Kind() == reflect.Struct {
			parentPrefix := getBsonTag(field)
			getSubStructFields(s.Field(i), parentPrefix+".", fields)
		}

		bsonTag := getBsonTag(field)
		if bsonTag != "" {
			fields[fieldValue.Addr().Interface()] = parentPrefix + bsonTag
		}
	}
}

func contains[T comparable](slice []T, item T) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (k *Kyte) Errors() []error {
	return k.errs
}

func (k *Kyte) setError(err error) {
	k.errs = append(k.errs, err)
}

func (k *Kyte) validateQueryFieldAndValue(field any, value any) (string, error) {
	if len(k.errs) > 0 {
		return "", k.errs[0]
	}

	if value == nil {
		return "", ErrNilValue
	}

	fieldType := reflect.TypeOf(field)
	if fieldType.Kind() != reflect.String && fieldType.Kind() != reflect.Ptr {
		return "", ErrFieldMustBePtrOrString
	}

	if fieldType.Kind() == reflect.String && field.(string) == "" {
		return "", ErrEmptyField
	}

	if fieldType.Kind() == reflect.Pointer && field == nil {
		return "", ErrNilPointerField
	}

	if k.fieldCheck {
		fieldName := ""
		if fieldType.Kind() == reflect.String {
			fieldName = field.(string)
		}

		ok := false
		if fieldType.Kind() == reflect.Ptr {
			_, ok = k.fields[field]
		}

		if !ok && !contains(k.fieldNames, fieldName) {
			return "", errors.Join(ErrNotValidFieldForQuery, errors.New(fmt.Sprintf("field: %s You can ignore this error by setting fieldCheck to false", fieldName)))
		}

	}

	fieldName, err := k.getFieldName(field)
	if err != nil {
		return "", err
	}

	return fieldName, nil
}

func (k *Kyte) hasErrors() bool {
	return len(k.errs) > 0
}

func (k *Kyte) getFieldName(field any) (string, error) {
	if reflect.TypeOf(field).Kind() == reflect.String {
		return field.(string), nil
	}

	if reflect.TypeOf(field).Kind() == reflect.Ptr {
		fieldName, ok := k.fields[field]
		if !ok {
			return "", ErrNotValidFieldForQuery
		}
		return fieldName, nil
	}

	return "", ErrFieldMustBePtrOrString
}
