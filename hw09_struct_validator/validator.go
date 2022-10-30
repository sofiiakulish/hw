package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const VALIDATE_TAG_NAME = "validate"

var NoValidationError = ValidationError{"", nil}

type ValidatorType string

const (
	MIN ValidatorType = "min"
	MAX ValidatorType = "max"
	LEN ValidatorType = "len"
	REGEXP ValidatorType = "regexp"
	IN ValidatorType = "in"
)

type ValidatorValue string

type Validator struct {
	Type  ValidatorType
	Value ValidatorValue
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var b strings.Builder

	for _, validationError := range v {
	    b.WriteString(fmt.Sprintf("Validation error for field %q: %s\n", validationError.Field, validationError.Err.Error()))
	}

	return b.String()
}

func Validate(v interface{}) error {
	var result ValidationErrors

    stT := reflect.TypeOf(v)
    stV := reflect.ValueOf(v)

    for i:=0; i<stT.NumField(); i++ {
        field := stT.Field(i)
        tag := field.Tag.Get(VALIDATE_TAG_NAME)
        if tag == "" {
            continue
        }

		validators, err := getValidators(tag)
		if err != nil {
		    return err
		}

		for _, validator := range validators {
			validationErrors, err := valid(field, reflect.Indirect(stV).FieldByName(field.Name).Interface(), validator)
			if err != nil {
			    return err
			}

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

func getValidators(tagValue string) ([]Validator, error) {
	var validator []Validator

	validatorStrings := strings.Split(tagValue, "|")

	for _, str := range validatorStrings {

		delimiterPos := strings.Index(str, ":")
		if delimiterPos == -1 {
			return nil, errors.New("Incorrect validator definition. Semicolon should be present")
		}

		typeString := str[0:delimiterPos]

		ruleString := str[delimiterPos+1:]
		if ruleString == "" {
			return nil, errors.New("Incorrect validator definition. Value should be present after semicolon")
		}

		validator = append(validator, Validator{ValidatorType(typeString), ValidatorValue(ruleString)})
	}

	if len(validator) == 0 {
		return nil, nil
	}

	return validator, nil
}

func valid(f reflect.StructField, fieldValue interface{}, validator Validator) ([]ValidationError, error) {
	switch f.Type.Kind() {
	case reflect.Slice, reflect.Array:
		var result []ValidationError

		array := reflect.ValueOf(fieldValue)

		for i := 0; i < array.Len(); i++ {
		    oneEntryResult, err := validOneEntry(f.Name, array.Index(i).Interface(), validator)
		    if err != nil {
		        return nil, err
		    }
			result = append(result, oneEntryResult...)
		}

		return result, nil
	}

	return validOneEntry(f.Name, fieldValue, validator)
}

func validOneEntry(fieldName string, fieldValue interface{}, validator Validator) ([]ValidationError, error) {
	switch validator.Type {
	case "min":
		value, typeErr := getInt(fieldValue)
		if typeErr != nil {
		    return nil, typeErr
		}
		validationErr, err := validateMin(fieldName, value, validator)
		if err != nil {
		    return nil, err
		}
		if validationErr != NoValidationError {
			return []ValidationError{validationErr}, nil
		}
	case "max":
		value, typeErr := getInt(fieldValue)
		if typeErr != nil {
		    return nil, typeErr
		}
		validationErr, err := validateMax(fieldName, value, validator)
        if err != nil {
            return nil, err
        }
		if validationErr != NoValidationError {
			return []ValidationError{validationErr}, nil
		}
	case "len":
		value, typeErr := getString(fieldValue)
		if typeErr != nil {
		    return nil, typeErr
		}
		validationErr, err := validateLen(fieldName, value, validator)
        if err != nil {
            return nil, err
        }
		if validationErr != NoValidationError {
			return []ValidationError{validationErr}, nil
		}
	case "regexp":
		value, typeErr := getString(fieldValue)
		if typeErr != nil {
		    return nil, typeErr
		}
		validationErr, err := validateRegexp(fieldName, value, validator)
        if err != nil {
            return nil, err
        }
		if validationErr != NoValidationError {
			return []ValidationError{validationErr}, nil
		}
	case "in":
		switch reflect.TypeOf(fieldValue).Kind() {
		case reflect.Int:
			value, typeErr := getInt(fieldValue)
            if typeErr != nil {
                return nil, typeErr
            }
			validationErr, err := validateInInt(fieldName, value, validator)
            if err != nil {
                return nil, err
            }
			if validationErr != NoValidationError {
				return []ValidationError{validationErr}, nil
			}
		case reflect.String:
			value, typeErr := getString(fieldValue)
            if typeErr != nil {
                return nil, typeErr
            }
			validationErr, err := validateInString(fieldName, value, validator)
            if err != nil {
                return nil, err
            }
			if validationErr != NoValidationError {
				return []ValidationError{validationErr}, nil
			}
		}
	default:
		return nil, errors.New("Unsupported validation type")
	}

	return nil, nil
}

func getInt(v interface{}) (int, error) {
	if reflect.TypeOf(v).Kind() != reflect.Int {
		return 0, fmt.Errorf("Value should have type int. %s provided", reflect.TypeOf(v).Name())
	}

	return int(reflect.ValueOf(v).Int()), nil
}

func getString(v interface{}) (string, error) {
	if reflect.TypeOf(v).Kind() != reflect.String {
		return "", fmt.Errorf("Value should have type string. %s provided", reflect.TypeOf(v).Kind())
	}

	return reflect.ValueOf(v).String(), nil
}

func validateMin(fieldName string, fieldValue int, validator Validator) (ValidationError, error) {
	validatorValue, err := strconv.Atoi(string(validator.Value))
	if err != nil {
		return NoValidationError, fmt.Errorf("Error getting int definition for the validator: %s", err)
	}

	if fieldValue < validatorValue {
		return ValidationError{fieldName, fmt.Errorf("Value should be more than %d. %d was provided", validatorValue, fieldValue)}, nil
	}

	return NoValidationError, nil
}

func validateMax(fieldName string, fieldValue int, validator Validator) (ValidationError, error) {
	validatorValue, err := strconv.Atoi(string(validator.Value))
	if err != nil {
		return NoValidationError, fmt.Errorf("Error getting int definition for the validator: %s", err)
	}

	if fieldValue > validatorValue {
		return ValidationError{fieldName, fmt.Errorf("Value should be less than %d. %d was provided", validatorValue, fieldValue)}, nil
	}

	return NoValidationError, nil
}

func validateInString(fieldName string, fieldValue string, validator Validator) (ValidationError, error) {
	validatorValues := strings.Split(string(validator.Value), ",")

	for _, val := range validatorValues {
		if val == fieldValue {
			return NoValidationError, nil
		}
	}

	return ValidationError{fieldName, fmt.Errorf("Value should be in the list %q. %q was provided", string(validator.Value), fieldValue)}, nil
}

func validateInInt(fieldName string, fieldValue int, validator Validator) (ValidationError, error) {
	validatorValuesString := strings.Split(string(validator.Value), ",")
	validatorValues := make([]int, len(validatorValuesString))

	for k, v := range validatorValuesString {
		value, err := strconv.Atoi(v)
		if err != nil {
			return NoValidationError, fmt.Errorf("Error parsing validation value: %s", err)
		}

		validatorValues[k] = value
	}

	for _, val := range validatorValues {
		if val == fieldValue {
			return NoValidationError, nil
		}
	}

	return ValidationError{fieldName, fmt.Errorf("Value should be in the list %q. %d was provided", string(validator.Value), fieldValue)}, nil
}

func validateLen(fieldName string, fieldValue string, validator Validator) (ValidationError, error) {
	validatorValue, _ := strconv.Atoi(string(validator.Value))

	if len(fieldValue) != validatorValue {
		return ValidationError{fieldName, fmt.Errorf("Value should have length %d. %d (%q) provided", validatorValue, len(fieldValue), fieldValue)}, nil
	}

	return NoValidationError, nil
}

func validateRegexp(fieldName string, fieldValue string, validator Validator) (ValidationError, error) {
	validatorValue := string(validator.Value)

	matched, err := regexp.MatchString(validatorValue, fieldValue)
	if err != nil {
		return NoValidationError, fmt.Errorf("Error regexp: %s", err)
	}

	if matched == false {
		return ValidationError{fieldName, fmt.Errorf("Value should match regexp %q. Value %q is not match", validatorValue, fieldValue)}, nil
	}

	return NoValidationError, nil
}
