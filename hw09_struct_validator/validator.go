package hw09structvalidator

import (
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const VALIDATE_TAG_NAME = "validate"

var NilValidationError = ValidationError{"", nil}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

type ValidatorType int

const (
	MIN ValidatorType = iota
	MAX
	LEN
	REGEXP
	IN
)

type ValidatorValue string

type Validator struct {
	Type  ValidatorType
	Value ValidatorValue
}

func (v ValidationErrors) Error() string {
	var result string

	for _, validationError := range v {
		result += fmt.Sprintf("Validation error for field %q: %s\n", validationError.Field, validationError.Err.Error())
	}

	return result
}

func Validate(v interface{}) error {
	var result ValidationErrors

	fields := structs.Fields(v)

	for _, field := range fields {
		tag := field.Tag(VALIDATE_TAG_NAME)
		if tag == "" {
			continue
		}

		validators := getValidators(tag)

		for _, validator := range validators {
			validationErrors := valid(field, validator)
			if validationErrors != nil {
				result = append(result, validationErrors...)
			}
		}
	}

	if len(result) == 0 {
		return nil
	}

	return errors.New(result.Error())
}

func getValidators(tagValue string) []Validator {
	var validator []Validator

	validatorStrings := strings.Split(tagValue, "|")

	for _, str := range validatorStrings {

		delimiterPos := strings.Index(str, ":")
		if delimiterPos == -1 {
			panic("Incorrect validator definition. Semicolon should be present")
		}

		typeString := str[0:delimiterPos]
		typeValue := getValidatorTypeByString(typeString)

		ruleString := str[delimiterPos+1:]
		if ruleString == "" {
			panic("Incorrect validator definition. Value should be present after semicolon")
		}
		ruleValue := getValidatorValueByString(ruleString)

		validator = append(validator, Validator{typeValue, ruleValue})
	}

	if len(validator) == 0 {
		return nil
	}

	return validator
}

func getValidatorTypeByString(typeString string) ValidatorType {
	switch typeString {
	case "min":
		return MIN
	case "max":
		return MAX
	case "len":
		return LEN
	case "regexp":
		return REGEXP
	case "in":
		return IN
	default:
		panic("Unknown validation type '" + typeString + "'")
	}
}

func getValidatorValueByString(ruleString string) ValidatorValue {
	return ValidatorValue(ruleString)
}

func valid(f *structs.Field, validator Validator) []ValidationError {
	rt := reflect.TypeOf(f.Value())
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		var result []ValidationError

		array := reflect.ValueOf(f.Value())

		for i := 0; i < array.Len(); i++ {
			result = append(result, validOneEntry(f.Name(), array.Index(i).Interface(), validator)...)
		}

		return result
	}

	return validOneEntry(f.Name(), f.Value(), validator)
}

func validOneEntry(fieldName string, fieldValue interface{}, validator Validator) []ValidationError {
	switch validator.Type {
	case MIN:
		value := getIntOrPanic(fieldValue)
		err := validateMin(fieldName, value, validator)
		if err != NilValidationError {
			return []ValidationError{err}
		}
	case MAX:
		value := getIntOrPanic(fieldValue)
		err := validateMax(fieldName, value, validator)
		if err != NilValidationError {
			return []ValidationError{err}
		}
	case LEN:
		value := getStringOrPanic(fieldValue)

		err := validateLen(fieldName, value, validator)
		if err != NilValidationError {
			return []ValidationError{err}
		}
	case REGEXP:
		value := getStringOrPanic(fieldValue)

		err := validateRegexp(fieldName, value, validator)
		if err != NilValidationError {
			return []ValidationError{err}
		}
	case IN:
		switch reflect.TypeOf(fieldValue).Kind() {
		case reflect.Int:
			value := getIntOrPanic(fieldValue)

			err := validateInInt(fieldName, value, validator)
			if err != NilValidationError {
				return []ValidationError{err}
			}
		case reflect.String:
			value := getStringOrPanic(fieldValue)

			err := validateInString(fieldName, value, validator)
			if err != NilValidationError {
				return []ValidationError{err}
			}
		}
	default:
		panic("Unsupported validation type")
	}

	return nil
}

func getIntOrPanic(v interface{}) int {
	if reflect.TypeOf(v).Kind() != reflect.Int {
		panic(fmt.Errorf("Value should have type int. %s provided", reflect.TypeOf(v).Name()))
	}

	return int(reflect.ValueOf(v).Int())
}

func getStringOrPanic(v interface{}) string {
	if reflect.TypeOf(v).Kind() != reflect.String {
		panic(fmt.Errorf("Value should have type string. %s provided", reflect.TypeOf(v).Kind()))
	}

	return reflect.ValueOf(v).String()
}

func validateMin(fieldName string, fieldValue int, validator Validator) ValidationError {
	validatorValue, err := strconv.Atoi(string(validator.Value))
	if err != nil {
		panic(fmt.Sprintf("Error geting int definition for the validator: %s", err))
	}

	if fieldValue < validatorValue {
		return ValidationError{fieldName, fmt.Errorf("Value should be more than %d. %d was provided", validatorValue, fieldValue)}
	}

	return NilValidationError
}

func validateMax(fieldName string, fieldValue int, validator Validator) ValidationError {
	validatorValue, err := strconv.Atoi(string(validator.Value))
	if err != nil {
		panic(err)
	}

	if fieldValue > validatorValue {
		return ValidationError{fieldName, fmt.Errorf("Value should be less than %d. %d was provided", validatorValue, fieldValue)}
	}

	return NilValidationError
}

func validateInString(fieldName string, fieldValue string, validator Validator) ValidationError {
	validatorValues := strings.Split(string(validator.Value), ",")

	for _, val := range validatorValues {
		if val == fieldValue {
			return NilValidationError
		}
	}

	return ValidationError{fieldName, fmt.Errorf("Value should be in the list %q. %q was provided", string(validator.Value), fieldValue)}
}

func validateInInt(fieldName string, fieldValue int, validator Validator) ValidationError {
	validatorValuesString := strings.Split(string(validator.Value), ",")
	validatorValues := make([]int, len(validatorValuesString))

	for k, v := range validatorValuesString {
		value, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Errorf("Error parsing validation value: %s", err))
		}

		validatorValues[k] = value
	}

	for _, val := range validatorValues {
		if val == fieldValue {
			return NilValidationError
		}
	}

	return ValidationError{fieldName, fmt.Errorf("Value should be in the list %q. %d was provided", string(validator.Value), fieldValue)}
}

func validateLen(fieldName string, fieldValue string, validator Validator) ValidationError {
	validatorValue, _ := strconv.Atoi(string(validator.Value))

	if len(fieldValue) != validatorValue {
		return ValidationError{fieldName, fmt.Errorf("Value should have length %d. %d (%q) provided", validatorValue, len(fieldValue), fieldValue)}
	}

	return NilValidationError
}

func validateRegexp(fieldName string, fieldValue string, validator Validator) ValidationError {
	validatorValue := string(validator.Value)

	matched, err := regexp.MatchString(validatorValue, fieldValue)
	if err != nil {
		panic(fmt.Errorf("Error regexp: %s", err))
	}

	if matched == false {
		return ValidationError{fieldName, fmt.Errorf("Value should match regexp %q. Value %q is not match", validatorValue, fieldValue)}
	}

	return NilValidationError
}
