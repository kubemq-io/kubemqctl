package config

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Validator func(val interface{}) error

func IsRequired() Validator {
	return func(val interface{}) error {
		// the reflect value of the result
		value := reflect.ValueOf(val)

		// if the value passed in is the zero value of the appropriate type
		if isZero(value) && value.Kind() != reflect.Bool {
			return errors.New("value is required")
		}
		return nil
	}
}

func IsUint() Validator {

	return func(val interface{}) error {
		if str, ok := val.(string); ok {
			// if the string is longer than the given value
			v, err := strconv.Atoi(str)
			if err != nil {
				return fmt.Errorf("value is not an integer")
			}
			if v < 0 {
				return fmt.Errorf("value must be greater then 0")
			}

		} else {
			// otherwise we cannot convert the value into a string and cannot enforce length
			return fmt.Errorf("cannot enforce length on response of type %v", reflect.TypeOf(val).Name())
		}

		// the input is fine
		return nil
	}
}

func IsInt() Validator {

	return func(val interface{}) error {
		if str, ok := val.(string); ok {
			// if the string is longer than the given value
			_, err := strconv.Atoi(str)
			if err != nil {
				return fmt.Errorf("value is not an integer")
			}

		} else {
			// otherwise we cannot convert the value into a string and cannot enforce length
			return fmt.Errorf("cannot enforce length on response of type %v", reflect.TypeOf(val).Name())
		}

		// the input is fine
		return nil
	}
}

func IsIntInRange(min, max int) Validator {

	return func(val interface{}) error {
		if str, ok := val.(string); ok {
			// if the string is longer than the given value
			v, err := strconv.Atoi(str)
			if err != nil {
				return fmt.Errorf("value is not an integer")
			}

			if v < min {
				return fmt.Errorf("value must be greater then %d", min)
			}
			if v > max {
				return fmt.Errorf("value must be lower then %d", max)
			}
		} else {
			// otherwise we cannot convert the value into a string and cannot enforce length
			return fmt.Errorf("cannot enforce length on response of type %v", reflect.TypeOf(val).Name())
		}

		// the input is fine
		return nil
	}
}

func IsValidHostPort() Validator {
	// return a validator that checks the length of the string
	return func(val interface{}) error {
		if str, ok := val.(string); ok {
			if str == "" {
				return nil
			}
			pair := strings.Split(str, ":")
			if len(pair) != 2 {
				return fmt.Errorf("value must be in a form host:port")
			}
			v, err := strconv.Atoi(pair[1])
			if err != nil {
				return fmt.Errorf("port is not an integer")
			}
			if v < 0 {
				return fmt.Errorf("port must be greater then 0")
			}
			if v > 65535 {
				return fmt.Errorf("port must be lower then 65535")
			}
		} else {
			// otherwise we cannot convert the value into a string and cannot enforce length
			return fmt.Errorf("cannot enforce length on response of type %v", reflect.TypeOf(val).Name())
		}

		// the input is fine
		return nil
	}
}

func IsMaxLength(length int) Validator {
	// return a validator that checks the length of the string
	return func(val interface{}) error {
		if str, ok := val.(string); ok {

			if len([]rune(str)) > length {

				return fmt.Errorf("value is too long. Max length is %v", length)
			}
		} else {
			// otherwise we cannot convert the value into a string and cannot enforce length
			return fmt.Errorf("cannot enforce length on response of type %v", reflect.TypeOf(val).Name())
		}

		// the input is fine
		return nil
	}
}

func IsExactLength(length int) Validator {
	// return a validator that checks the length of the string
	return func(val interface{}) error {
		if str, ok := val.(string); ok {

			if len([]rune(str)) != length {
				return fmt.Errorf("value length must be  %v", length)
			}
		} else {
			// otherwise we cannot convert the value into a string and cannot enforce length
			return fmt.Errorf("cannot enforce length on response of type %v", reflect.TypeOf(val).Name())
		}

		// the input is fine
		return nil
	}
}

func IsValidToken() Validator {
	// return a validator that checks the length of the string
	return func(val interface{}) error {
		if str, ok := val.(string); ok {
			if str == "" {
				return nil
			}
			segments := strings.Split(str, "-")
			if len(segments) != 5 {
				return fmt.Errorf("value is not valid token")
			}
			if len(segments[0]) != 8 ||
				len(segments[1]) != 4 ||
				len(segments[2]) != 4 ||
				len(segments[3]) != 4 ||
				len(segments[4]) != 12 {
				return fmt.Errorf("value is not valid token")
			}

		} else {
			// otherwise we cannot convert the value into a string and cannot enforce length
			return fmt.Errorf("cannot enforce length on response of type %v", reflect.TypeOf(val).Name())
		}

		// the input is fine
		return nil
	}
}

func IsMinLength(length int) Validator {
	// return a validator that checks the length of the string
	return func(val interface{}) error {
		if str, ok := val.(string); ok {
			// if the string is shorter than the given value
			if len([]rune(str)) < length {
				// yell loudly
				return fmt.Errorf("value is too short. Min length is %v", length)
			}
		} else {
			// otherwise we cannot convert the value into a string and cannot enforce length
			return fmt.Errorf("cannot enforce length on response of type %v", reflect.TypeOf(val).Name())
		}

		// the input is fine
		return nil
	}
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	}

	// compare the types directly with more general coverage
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
