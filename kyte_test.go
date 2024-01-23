package kyte

import (
	"errors"
	"testing"
)

type TestTodo struct {
	ID      string `bson:"id"`
	Name    string `bson:"name"`
	Message string `bson:"message"`
}

type TestArrayStruct struct {
	UserID   string     `bson:"user_id"`
	Username string     `bson:"username"`
	Type     string     `bson:"type"`
	Todos    []TestTodo `bson:"todos"`
}

type TestArrayStructWithPointer struct {
	UserID   string      `bson:"user_id"`
	Username string      `bson:"username"`
	Type     string      `bson:"type"`
	Todos    *[]TestTodo `bson:"todos"`
}

type TestStructWithPointerArray struct {
	UserID   string      `bson:"user_id"`
	Username string      `bson:"username"`
	Type     string      `bson:"type"`
	Todos    []*TestTodo `bson:"todos"`
}

type TestWithNestedStruct struct {
	UserID   string   `bson:"user_id"`
	Username string   `bson:"username"`
	Type     string   `bson:"type"`
	Todo     TestTodo `bson:"todo"`
}

type TestWithNestedStructWithPointer struct {
	UserID   string    `bson:"user_id"`
	Username string    `bson:"username"`
	Type     string    `bson:"type"`
	Todo     *TestTodo `bson:"todo"`
}

func testKyteFieldAndSource(t *testing.T, kyte *kyte, fields map[any]string, errCount int, fieldCount int) {
	if len(kyte.fields) != fieldCount {
		t.Errorf("kyte.fields should be %v but got %v", fieldCount, len(kyte.fields))
	}

	if len(kyte.fieldNames) != fieldCount {
		t.Errorf("kyte.fieldNames should be %v but got %v", fieldCount, len(kyte.fieldNames))
	}

	if len(kyte.errs) != errCount {
		t.Errorf("kyte.errs should be empty slice but got %v", kyte.errs)
	}

	for ptr, field := range fields {
		if _, ok := kyte.fields[ptr]; !ok {
			t.Errorf("kyte.fields should have %v field", field)
		}
	}
}

func TestNewKyte(t *testing.T) {
	t.Parallel()

	t.Run("with source nil should return empty kyte", func(t *testing.T) {
		kyte := newKyte(nil, true)
		if kyte.source != nil {
			t.Errorf("kyte.source should be nil but got %v", kyte.source)
		}
		testKyteFieldAndSource(t, kyte, nil, 0, 0)
	})

	t.Run("with source not nil should return kyte with fields", func(t *testing.T) {
		t.Run("zero value", func(t *testing.T) {
			source := &TestArrayStruct{}
			fields := map[any]string{
				&source.UserID:   "user_id",
				&source.Username: "username",
				&source.Type:     "type",
				&source.Todos:    "todos",
			}
			kyte := newKyte(source, true)
			if kyte.source == nil {
				t.Errorf("kyte.source should not be nil")
			}
			testKyteFieldAndSource(t, kyte, fields, 0, 7)
		})

		t.Run("non zero value", func(t *testing.T) {
			source := &TestArrayStruct{
				UserID:   "user_id",
				Username: "username",
				Type:     "type",
				Todos:    []TestTodo{},
			}
			fields := map[any]string{
				&source.UserID:   "user_id",
				&source.Username: "username",
				&source.Type:     "type",
				&source.Todos:    "todos",
			}
			kyte := newKyte(source, true)
			if kyte.source == nil {
				t.Errorf("kyte.source should not be nil")
			}
			testKyteFieldAndSource(t, kyte, fields, 0, 7)
		})
	})

	t.Run("with source with pointer of struct should return kyte with fields", func(t *testing.T) {
		t.Run("zero value", func(t *testing.T) {
			source := &TestArrayStructWithPointer{}
			fields := map[any]string{
				&source.UserID:   "user_id",
				&source.Username: "username",
				&source.Type:     "type",
				&source.Todos:    "todos",
			}
			kyte := newKyte(source, true)
			if kyte.source == nil {
				t.Errorf("kyte.source should not be nil")
			}

			testKyteFieldAndSource(t, kyte, fields, 0, 7)
		})
		t.Run("non zero value", func(t *testing.T) {

			source := &TestArrayStructWithPointer{Todos: &[]TestTodo{}}
			fields := map[any]string{
				&source.UserID:   "user_id",
				&source.Username: "username",
				&source.Type:     "type",
				&source.Todos:    "todos",
			}
			kyte := newKyte(source, true)
			if kyte.source == nil {
				t.Errorf("kyte.source should not be nil")
			}

			testKyteFieldAndSource(t, kyte, fields, 0, 7)
		})
	})

	t.Run("with source with array of pointer of struct should return kyte with fields", func(t *testing.T) {
		t.Run("zero value", func(t *testing.T) {
			source := &TestStructWithPointerArray{}
			fields := map[any]string{
				&source.UserID: "user_id",
				&source.Todos:  "todos",
				&source.Type:   "type",
			}

			kyte := newKyte(source, true)
			if kyte.source == nil {
				t.Errorf("kyte.source should not be nil")
			}
			testKyteFieldAndSource(t, kyte, fields, 0, 7)
		})
		t.Run("non zero value", func(t *testing.T) {
			source := &TestStructWithPointerArray{Todos: []*TestTodo{}}
			fields := map[any]string{
				&source.UserID: "user_id",
				&source.Todos:  "todos",
				&source.Type:   "type",
			}

			kyte := newKyte(source, true)
			if kyte.source == nil {
				t.Errorf("kyte.source should not be nil")
			}
			testKyteFieldAndSource(t, kyte, fields, 0, 7)
		})
	})

	t.Run("with source with nested struct should return kyte with fields", func(t *testing.T) {

		t.Run("zero value", func(t *testing.T) {
			source := &TestWithNestedStruct{}
			fields := map[any]string{
				&source.UserID:       "user_id",
				&source.Username:     "username",
				&source.Type:         "type",
				&source.Todo:         "todo",
				&source.Todo.ID:      "todo.id",
				&source.Todo.Name:    "todo.name",
				&source.Todo.Message: "todo.message",
			}

			kyte := newKyte(source, true)
			if kyte.source == nil {
				t.Errorf("kyte.source should not be nil")
			}
			testKyteFieldAndSource(t, kyte, fields, 0, 7)
		})
		t.Run("non zero value", func(t *testing.T) {
			source := &TestWithNestedStruct{Todo: TestTodo{}}
			fields := map[any]string{
				&source.UserID:       "user_id",
				&source.Username:     "username",
				&source.Type:         "type",
				&source.Todo:         "todo",
				&source.Todo.ID:      "todo.id",
				&source.Todo.Name:    "todo.name",
				&source.Todo.Message: "todo.message",
			}

			kyte := newKyte(source, true)
			if kyte.source == nil {
				t.Errorf("kyte.source should not be nil")
			}
			testKyteFieldAndSource(t, kyte, fields, 0, 7)
		})
	})

	t.Run("with source with nested struct with pointer should return kyte with fields", func(t *testing.T) {
		source := &TestWithNestedStructWithPointer{Todo: &TestTodo{}}
		fields := map[any]string{
			&source.UserID:       "user_id",
			&source.Username:     "username",
			&source.Type:         "type",
			&source.Todo:         "todo",
			&source.Todo.ID:      "todo.id",
			&source.Todo.Name:    "todo.name",
			&source.Todo.Message: "todo.message",
		}

		kyte := newKyte(source, true)
		testKyteFieldAndSource(t, kyte, fields, 0, 7)
	})

	t.Run("not pointer source should return kyte with error", func(t *testing.T) {
		source := TestArrayStruct{}
		kyte := newKyte(source, true)
		if !kyte.hasErrors() {
			t.Errorf("kyte should have errors")
		}

		if kyte.errs[0] != ErrNotPtrSource {
			t.Errorf("kyte should have error %v but got %v", ErrNotPtrSource, kyte.errs[0])
		}
	})

	t.Run("source unexported field should not be added to kyte", func(t *testing.T) {
		type TestAnonymousStruct struct {
			Name string `bson:"name"`
			age  int    `bson:"age"`
		}

		source := &TestAnonymousStruct{}
		fields := map[any]string{
			&source.Name: "name",
		}

		kyte := newKyte(source, true)
		testKyteFieldAndSource(t, kyte, fields, 0, 1)
	})

	t.Run("array of pointer of struct non zero value ", func(t *testing.T) {
		type TestAnonymousWithSlice struct {
			Name  string      `bson:"name"`
			Todos []*TestTodo `bson:"todos"`
		}
		todo := TestTodo{}
		source := &TestAnonymousWithSlice{Todos: []*TestTodo{
			&todo,
		}}
		fields := map[any]string{
			&source.Name:  "name",
			&source.Todos: "todos",
			&todo.ID:      "todos.id",
			&todo.Name:    "todos.name",
			&todo.Message: "todos.message",
		}

		kyte := newKyte(source, true)
		testKyteFieldAndSource(t, kyte, fields, 0, 5)
	})

	t.Run("array of pointer of struct zero value ", func(t *testing.T) {
		type TestAnonymousWithSlice struct {
			Name  string      `bson:"name"`
			Todos []*TestTodo `bson:"todos"`
		}
		source := &TestAnonymousWithSlice{}
		fields := map[any]string{
			&source.Name:  "name",
			&source.Todos: "todos",
		}

		kyte := newKyte(source, true)
		testKyteFieldAndSource(t, kyte, fields, 0, 5)
	})

	t.Run("pointer array of pointer of struct non zero value ", func(t *testing.T) {
		type TestAnonymousWithSlicePointer struct {
			Name  string       `bson:"name"`
			Todos *[]*TestTodo `bson:"todos"`
		}

		todo := TestTodo{}
		source := &TestAnonymousWithSlicePointer{Todos: &[]*TestTodo{
			&todo,
		}}
		fields := map[any]string{
			&source.Name:  "name",
			&source.Todos: "todos",
			&todo.ID:      "todos.id",
			&todo.Name:    "todos.name",
			&todo.Message: "todos.message",
		}

		kyte := newKyte(source, true)
		testKyteFieldAndSource(t, kyte, fields, 0, 5)
	})

	t.Run("pointer array of pointer of struct zero value ", func(t *testing.T) {
		type TestAnonymousWithSlicePointer struct {
			Name  string       `bson:"name"`
			Todos *[]*TestTodo `bson:"todos"`
		}

		source := &TestAnonymousWithSlicePointer{}
		fields := map[any]string{
			&source.Name:  "name",
			&source.Todos: "todos",
		}

		kyte := newKyte(source, true)
		testKyteFieldAndSource(t, kyte, fields, 0, 5)
	})

	t.Run("ignore field if does not have bson tag", func(t *testing.T) {
		type TestAnonymousWithSlicePointer struct {
			Name string    `bson:"name"`
			Todo *TestTodo `bson:"todo"`
			Age  int       `bson:"-"`
			Time string
		}

		todo := TestTodo{}
		source := &TestAnonymousWithSlicePointer{
			Todo: &todo,
		}
		fields := map[any]string{
			&source.Name: "name",
			&source.Todo: "todo",
			&todo.ID:     "todo.id",
			&todo.Name:   "todo.name",
		}

		kyte := newKyte(source, true)
		testKyteFieldAndSource(t, kyte, fields, 0, 5)
	})
}

func TestKyteErros(t *testing.T) {
	t.Parallel()
	t.Run("not nil errors", func(t *testing.T) {
		kyte := newKyte("str", true)
		if kyte.Errors() == nil {
			t.Errorf("kyte.Errors() should not be nil but got %v", kyte.Errors())
		}
	})

	t.Run("nil errors", func(t *testing.T) {
		kyte := newKyte(nil, true)
		if kyte.Errors() != nil {
			t.Errorf("kyte.Errors() should be nil but got %v", kyte.Errors())
		}
	})
}

func TestKyteHasErrors(t *testing.T) {
	t.Parallel()
	t.Run("with errors", func(t *testing.T) {
		kyte := newKyte("str", true)
		if !kyte.hasErrors() {
			t.Errorf("kyte.hasErrors() should return true")
		}
	})

	t.Run("without errors", func(t *testing.T) {
		kyte := newKyte(nil, true)
		if kyte.hasErrors() {
			t.Errorf("kyte.hasErrors() should return false")
		}
	})
}

func TestKyteSetError(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		kyte := newKyte(nil, true)
		kyte.setError(ErrNotPtrSource)
		if kyte.Errors() == nil {
			t.Errorf("kyte.Errors() should not be nil but got %v", kyte.Errors())
		}
	})
}

func TestKyteValidate(t *testing.T) {
	t.Parallel()
	t.Run("with errors", func(t *testing.T) {
		kyte := newKyte("str", true)
		err := kyte.validate(&operation{
			field: "field",
			value: "value"},
		)
		if err == nil {
			t.Errorf("kyte.validate() should return error")
		}
	})

	t.Run("field is required but nil", func(t *testing.T) {
		kyte := newKyte(nil, true)
		err := kyte.validate(&operation{
			field:           nil,
			value:           "value",
			isFieldRequired: true,
		})
		if err != ErrNilField {
			t.Errorf("kyte.validate() should return error %v but got %v", ErrNilField, err)
		}
	})

	t.Run("value is required but nil", func(t *testing.T) {
		kyte := newKyte(nil, true)
		err := kyte.validate(&operation{
			field:           "field",
			value:           nil,
			isValueRequired: true,
		})
		if err != ErrNilValue {
			t.Errorf("kyte.validate() should return error %v but got %v", ErrNilValue, err)
		}
	})

	t.Run("field is required but not string or pointer", func(t *testing.T) {
		kyte := newKyte(nil, true)
		err := kyte.validate(&operation{
			field:           1,
			value:           "value",
			isFieldRequired: true,
		})
		if err != ErrFieldMustBePtrOrString {
			t.Errorf("kyte.validate() should return error %v but got %v", ErrFieldMustBePtrOrString, err)
		}
	})

	t.Run("field is required but string and empty", func(t *testing.T) {
		kyte := newKyte(nil, true)
		err := kyte.validate(&operation{
			field:           "",
			value:           "value",
			isFieldRequired: true,
		})
		if err != ErrEmptyField {
			t.Errorf("kyte.validate() should return error %v but got %v", ErrEmptyField, err)
		}
	})

	t.Run("field is required and check field and field valid", func(t *testing.T) {
		kyte := newKyte(&TestTodo{}, true)
		err := kyte.validate(&operation{
			field:           "name",
			value:           "value",
			isFieldRequired: true,
		})
		if err != nil {
			t.Errorf("kyte.validate() should not return error but got %v", err)
		}
	})

	t.Run("field is required and check field and field invalid", func(t *testing.T) {
		kyte := newKyte(&TestTodo{}, true)
		err := kyte.validate(&operation{
			field:           "invalid",
			value:           "value",
			isFieldRequired: true,
		})
		if !errors.Is(err, ErrNotValidFieldForQuery) {
			t.Errorf("kyte.validate() should return error %v but got %v", ErrNotValidFieldForQuery, err)
		}
	})
}

func TestKyteIsFieldValid(t *testing.T) {
	t.Parallel()

	t.Run("with errors", func(t *testing.T) {
		kyte := newKyte("str", true)
		err := kyte.isFieldValid("field")
		if err == nil {
			t.Errorf("kyte.isFieldValid() should return error")
		}
	})

	t.Run("string field and valid", func(t *testing.T) {
		kyte := newKyte(&TestTodo{}, true)

		err := kyte.isFieldValid("name")
		if err != nil {
			t.Errorf("kyte.isFieldValid() should not return error but got %v", err)
		}
	})

	t.Run("string field and invalid", func(t *testing.T) {
		kyte := newKyte(&TestTodo{}, true)

		err := kyte.isFieldValid("invalid")
		if !errors.Is(err, ErrNotValidFieldForQuery) {
			t.Errorf("kyte.isFieldValid() should return error %v but got %v", ErrNotValidFieldForQuery, err)
		}
	})

	t.Run("pointer field and valid", func(t *testing.T) {
		todo := &TestTodo{}
		kyte := newKyte(todo, true)

		err := kyte.isFieldValid(&todo.Message)
		if err != nil {
			t.Errorf("kyte.isFieldValid() should not return error but got %v", err)
		}
	})
}

func TestKyteGetFieldName(t *testing.T) {
	t.Parallel()

	t.Run("string field", func(t *testing.T) {
		kyte := newKyte(&TestTodo{}, true)
		field, err := kyte.getFieldName("name")
		if err != nil {
			t.Errorf("kyte.getFieldName() should not return error but got %v", err)
		}

		if field != "name" {
			t.Errorf("kyte.getFieldName() should return name but got %v", field)
		}
	})

	t.Run("pointer field", func(t *testing.T) {
		todo := &TestTodo{}
		kyte := newKyte(todo, true)
		field, err := kyte.getFieldName(&todo.Message)
		if err != nil {
			t.Errorf("kyte.getFieldName() should not return error but got %v", err)
		}

		if field != "message" {
			t.Errorf("kyte.getFieldName() should return message but got %v", field)
		}
	})

	t.Run("pointer field and not in fields", func(t *testing.T) {
		type temp struct {
			Temp string `bson:"temp"`
		}
		tempS := &temp{}
		todo := &TestTodo{}
		kyte := newKyte(todo, true)
		_, err := kyte.getFieldName(&tempS.Temp)
		if !errors.Is(err, ErrNotValidFieldForQuery) {
			t.Errorf("kyte.getFieldName() should return error %v but got %v", ErrNotValidFieldForQuery, err)
		}
	})

	t.Run("pointer field and field is nil", func(t *testing.T) {
		todo := &TestTodo{}
		kyte := newKyte(nil, false)
		_, err := kyte.getFieldName(&todo.ID)
		if err != ErrFieldMustBeString {
			t.Errorf("kyte.getFieldName() should return error %v but got %v", ErrFieldMustBeString, err)
		}
	})

	t.Run("field is not string or pointer", func(t *testing.T) {
		kyte := newKyte(nil, false)
		_, err := kyte.getFieldName(1)
		if err != ErrFieldMustBePtrOrString {
			t.Errorf("kyte.getFieldName() should return error %v but got %v", ErrFieldMustBePtrOrString, err)
		}
	})
}
