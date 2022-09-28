package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	InvalidLenIntDefinition struct {
		ID int `validate:"len:36"`
	}

	InvalidMinStringDefinition struct {
		Name string `validate:"min:36"`
	}

	InvalidMaxStringDefinition struct {
		Name string `validate:"max:36"`
	}

	InvalidRegexpIntDefinition struct {
		Name int `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	}
	InvalidLenDefinition struct {
		ID int `validate:"len:a"`
	}

	InvalidMinDefinition struct {
		ID int `validate:"min:a"`
	}

	InvalidMaxDefinition struct {
		ID int `validate:"max:a"`
	}
)

func TestValidateWithValidationErrors(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{in: User{"123", "Name", 20, "1@gmail.com", "admin", []string{"12345678901"}, []byte(`{"data": "test"}`)}, expectedErr: errors.New("Validation error for field \"ID\": Value should have length 36. 3 (\"123\") provided\n")},
		{in: User{"123", "Name", 12, "1@gmail.com", "admin", []string{"12345678901"}, []byte(`{"data": "test"}`)}, expectedErr: errors.New("Validation error for field \"ID\": Value should have length 36. 3 (\"123\") provided\nValidation error for field \"Age\": Value should be more than 18. 12 was provided\n")},
		{in: User{"123456789012345678901234567890123456", "Name", 20, "1gmail.com", "admin", []string{"12345678901"}, []byte(`{"data": "test"}`)}, expectedErr: errors.New(fmt.Sprintf("Validation error for field \"Email\": Value should match regexp %q. Value \"1gmail.com\" is not match\n", "^\\w+@\\w+\\.\\w+$"))},
		{in: User{"123456789012345678901234567890123456", "Name", 20, "1@gmail.com", "admin2", []string{"12345678901"}, []byte(`{"data": "test"}`)}, expectedErr: errors.New("Validation error for field \"Role\": Value should be in the list \"admin,stuff\". \"admin2\" was provided\n")},
		{in: User{"123456789012345678901234567890123456", "Name", 20, "1@gmail.com", "admin", []string{"123", "1234", "12345678901"}, []byte(`{"data": "test"}`)}, expectedErr: errors.New("Validation error for field \"Phones\": Value should have length 11. 3 (\"123\") provided\nValidation error for field \"Phones\": Value should have length 11. 4 (\"1234\") provided\n")},

		{in: App{"1234"}, expectedErr: errors.New("Validation error for field \"Version\": Value should have length 5. 4 (\"1234\") provided\n")},

		{in: Response{401, "Not Auhorized"}, expectedErr: errors.New("Validation error for field \"Code\": Value should be in the list \"200,404,500\". 401 was provided\n")},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Error(t, err)
			require.Equal(t, tt.expectedErr, err)
			_ = tt
		})
	}
}

func TestValidateSuccess(t *testing.T) {
	tests := []struct {
		in interface{}
	}{
		{in: User{"123456789012345678901234567890123456", "Name", 20, "1@gmail.com", "admin", []string{"12345678901"}, []byte(`{"data": "test"}`)}},
		{in: App{"12345"}},
		{in: Token{[]byte(`"test"`), []byte(`"test"`), []byte(`"test"`)}},
		{in: Response{200, "{}"}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.NoError(t, err)
			_ = tt
		})
	}
}

func TestInvalidDefinition(t *testing.T) {
	tests := []struct {
		in interface{}
	}{
		{in: InvalidLenIntDefinition{12345}},
		{in: InvalidMinStringDefinition{"Name"}},
		{in: InvalidMaxStringDefinition{"Name"}},
		{in: InvalidRegexpIntDefinition{12}},
		{in: InvalidLenDefinition{12345}},
		{in: InvalidMinDefinition{123}},
		{in: InvalidMaxDefinition{123}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("The code did not panic")
				}
			}()

			tt := tt
			t.Parallel()

			Validate(tt.in)
			_ = tt
		})
	}
}
