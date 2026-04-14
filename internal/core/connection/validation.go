package connection

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (ve ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "validation error"
	}
	msgs := make([]string, len(ve.Errors))
	for i, e := range ve.Errors {
		msgs[i] = fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return strings.Join(msgs, "; ")
}

func GetValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		jsonTag := fld.Tag.Get("json")
		// Extract just the field name, ignoring modifiers like "omitempty"
		// e.g., "region,omitempty" -> "region"
		if jsonTag == "" {
			return ""
		}
		parts := strings.Split(jsonTag, ",")
		return parts[0]
	})

	return v
}

func UnmarshalConfig[T any](configMap map[string]interface{}) (T, error) {
	var config T
	data, err := json.Marshal(configMap)
	if err != nil {
		return config, fmt.Errorf("marshaling config: %w", err)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("unmarshaling config: %w", err)
	}
	return config, nil
}

func ValidateConfig(connConfig interface{}) error {
	validate := GetValidator()
	errs := validate.Struct(connConfig)
	if errs == nil {
		return nil
	}

	validationErrs, ok := errs.(validator.ValidationErrors)
	if !ok {
		return errs
	}

	var errors []ValidationError
	for _, e := range validationErrs {
		errors = append(errors, ValidationError{
			Field:   getJSONFieldName(e),
			Message: getErrorMessage(e),
		})
	}

	return ValidationErrors{Errors: errors}
}

func ValidateConnectionConfig(connType string, configMap map[string]interface{}) error {
	source, err := GetRegistry().Get(connType)
	if err != nil {
		return fmt.Errorf("invalid connection type: %s", connType)
	}
	return source.Validate(configMap)
}

func getJSONFieldName(e validator.FieldError) string {
	namespace := e.Namespace()
	parts := strings.SplitN(namespace, ".", 2)
	if len(parts) > 1 {
		return parts[1]
	}
	return e.Field()
}

func getErrorMessage(e validator.FieldError) string {
	field := getJSONFieldName(e)

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "required_if":
		return fmt.Sprintf("%s is required", field)
	case "required_with":
		return fmt.Sprintf("%s is required when %s is specified", field, e.Param())
	case "min":
		if e.Type().Kind() == reflect.String {
			return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
		}
		return fmt.Sprintf("%s must be at least %s", field, e.Param())
	case "max":
		if e.Type().Kind() == reflect.String {
			return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
		}
		return fmt.Sprintf("%s must be at most %s", field, e.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, e.Param())
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, e.Param())
	case "hostname":
		return fmt.Sprintf("%s must be a valid hostname", field)
	case "hostname_port":
		return fmt.Sprintf("%s must be a valid hostname:port", field)
	case "ip":
		return fmt.Sprintf("%s must be a valid IP address", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
