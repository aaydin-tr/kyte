package kyte

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
	ErrNotStruct    = errors.New("source is not a pointer of a struct")

	ErrEmptyField             = errors.New("field is empty use string or pointer of an source struct field")
	ErrNilPointerField        = errors.New("field is nil pointer")
	ErrFieldMustBePtrOrString = errors.New("field must be string or pointer of an source struct field")

	ErrNotValidFieldForQuery = errors.New("this field is not in the source struct")

	ErrNilField          = errors.New("field is nil")
	ErrFieldMustBeString = errors.New("field must be string")

	ErrInvalidBsonType = errors.New("invalid bson type")

	ErrValueMustBeSlice = errors.New("value must be slice")
	ErrRegexCannotBeNil = errors.New("regex cannot be nil")
)

const (
	// Alias for mongo _id field
	UnderScoreID = "_id"
	// Alias for mongo _id field with dollar sign prefix
	UnderScoreIDWithDollar = "$_id"
)

type Options struct {
	// Source is the struct that will be used to check if the field is valid for query based on the struct bson tags.
	source any

	// ValidateField is true by default, it will check if the field is valid for query based on the source struct bson tags.
	//
	// Default: true
	validateField bool
}

type OptionFunc func(*Options)

func ValidateField(validateField bool) OptionFunc {
	return func(o *Options) {
		o.validateField = validateField
	}
}

func Source(source any) OptionFunc {
	return func(o *Options) {
		o.source = source
	}
}

type kyte struct {
	source     any
	fields     map[any]string
	fieldNames []string
	err        error
	checkField bool
}

func newKyte(source any, checkField bool) *kyte {
	if source == nil {
		return &kyte{}
	}
	kyte := &kyte{fields: make(map[any]string), checkField: checkField}
	kyte.setSourceAndPrepareFields(source)
	return kyte
}

// TODO refactor this function
func (k *kyte) setSourceAndPrepareFields(source any) {
	k.source = source
	k.fields = make(map[any]string)
	k.fieldNames = []string{}

	if reflect.ValueOf(source).Kind() != reflect.Ptr {
		k.err = ErrNotPtrSource
		return
	}

	if reflect.ValueOf(source).Kind() == reflect.Ptr && reflect.ValueOf(source).Elem().Kind() != reflect.Struct {
		k.err = ErrNotStruct
		return
	}

	k.source = source
	v := reflect.ValueOf(source).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldValue := v.Field(i)
		kind := fieldValue.Kind()

		if !field.IsExported() {
			continue
		}

		bsonTag := getBsonTag(field)
		if bsonTag == "" {
			continue
		}

		var addr any
		if fieldValue.CanAddr() && fieldValue.Addr().Interface() != nil {
			addr = fieldValue.Addr().Interface()
		}

		if kind == reflect.Slice {
			sliceType := fieldValue.Type().Elem()
			if sliceType.Kind() == reflect.Ptr {
				sliceType = sliceType.Elem()
			}

			if sliceType.Kind() == reflect.Struct {
				if fieldValue.Len() > 0 {
					firstElem := fieldValue.Index(0)
					firstElemType := firstElem.Type()
					if firstElemType.Kind() == reflect.Ptr {
						firstElem = firstElem.Elem()
					}

					getSubStructFields(firstElem, bsonTag, k.fields)
				} else {
					sliceElem := reflect.New(sliceType).Elem()
					getSubStructFields(sliceElem, bsonTag, k.fields)
				}
			}

		}

		if kind == reflect.Ptr {
			underlyingValue := reflect.New(fieldValue.Type().Elem()).Elem()
			underlyingKind := underlyingValue.Kind()
			if underlyingKind == reflect.Struct {
				getSubStructFields(fieldValue.Elem(), bsonTag, k.fields)
			}

			if underlyingKind == reflect.Slice {
				underlyingValue = reflect.New(underlyingValue.Type().Elem()).Elem()
				underlyingType := underlyingValue.Type()
				underlyingKind = underlyingValue.Kind()
				if underlyingKind == reflect.Ptr {
					underlyingType = underlyingType.Elem()
				}

				if underlyingType.Kind() == reflect.Struct {
					if !fieldValue.IsZero() && fieldValue.Elem().Len() > 0 {
						firstElem := fieldValue.Elem().Index(0)
						firstElemType := firstElem.Type()
						if firstElemType.Kind() == reflect.Ptr {
							firstElem = firstElem.Elem()
						}
						getSubStructFields(firstElem, bsonTag, k.fields)
					} else {
						sliceElem := reflect.New(underlyingType).Elem()
						getSubStructFields(sliceElem, bsonTag, k.fields)
					}
				}

			}

		}

		if field.Type.Kind() == reflect.Struct {
			getSubStructFields(v.Field(i), bsonTag, k.fields)
		}

		if bsonTag != "" && addr != nil {
			k.fields[addr] = bsonTag
		}
	}

	for _, v := range k.fields {
		k.fieldNames = append(k.fieldNames, v)
	}
}

func (k *kyte) setError(err error) {
	if k.err == nil {
		k.err = err
	}
}

type operation struct {
	operator        string
	field           any
	value           any
	isFieldRequired bool
}

func (k *kyte) validate(opt *operation) error {
	if k.hasErrors() {
		return k.err
	}

	if opt.isFieldRequired && opt.field == nil {
		return ErrNilField
	}

	if opt.isFieldRequired {
		fieldType := reflect.TypeOf(opt.field)
		if fieldType.Kind() != reflect.String && fieldType.Kind() != reflect.Ptr {
			return ErrFieldMustBePtrOrString
		}

		if fieldType.Kind() == reflect.String && opt.field.(string) == "" {
			return ErrEmptyField
		}
	}

	if opt.isFieldRequired && (k.checkField && k.fields != nil) {
		if err := k.isFieldValid(opt.field); err != nil {
			return err
		}
	}

	return nil
}

func (k *kyte) isFieldValid(field any) error {
	if k.hasErrors() {
		return k.err
	}

	fieldType := reflect.TypeOf(field)

	fieldName := ""
	if fieldType.Kind() == reflect.String {
		fieldName = field.(string)
	}

	ok := false
	if fieldType.Kind() == reflect.Ptr {
		_, ok = k.fields[field]
	}

	if !ok && !contains(k.fieldNames, fieldName) {
		return errors.Join(ErrNotValidFieldForQuery, fmt.Errorf("field: %s You can ignore this error by setting checkField to false", fieldName))
	}

	return nil
}

func (k *kyte) hasErrors() bool {
	return k.err != nil
}

func (k *kyte) getFieldName(field any) (string, error) {
	if reflect.TypeOf(field).Kind() == reflect.String {
		return field.(string), nil
	}

	if reflect.TypeOf(field).Kind() == reflect.Ptr && k.fields == nil {
		return "", ErrFieldMustBeString
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

func getBsonTag(field reflect.StructField) string {
	bsonTag := field.Tag.Get("bson")
	if bsonTag == "" || bsonTag == "-" {
		return ""
	}

	splitBsonTag := strings.Split(bsonTag, ",")

	for _, tag := range splitBsonTag {
		if !contains(builtinTags, tag) || tag == "-" {
			return tag
		}
	}

	return ""
}

func getSubStructFields(s reflect.Value, parentPrefix string, fields map[any]string) {
	if parentPrefix != "" {
		parentPrefix += "."
	}

	for i := 0; i < s.NumField(); i++ {
		field := s.Type().Field(i)
		fieldValue := s.Field(i)
		if !field.IsExported() {
			continue
		}

		if field.Type.Kind() == reflect.Struct {
			parentPrefix := getBsonTag(field)
			getSubStructFields(s.Field(i), parentPrefix+".", fields)
		}

		bsonTag := getBsonTag(field)
		if bsonTag != "" && fieldValue.CanAddr() {
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
