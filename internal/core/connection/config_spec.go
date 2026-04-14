package connection

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type FieldType string

const (
	FieldTypeString      FieldType = "string"
	FieldTypeInt         FieldType = "int"
	FieldTypeBool        FieldType = "bool"
	FieldTypeSelect      FieldType = "select"
	FieldTypeMultiselect FieldType = "multiselect"
	FieldTypePassword    FieldType = "password"
	FieldTypeObject      FieldType = "object"
)

// ShowWhen defines conditional field visibility
type ShowWhen struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type ConfigField struct {
	Name        string        `json:"name"`
	Type        FieldType     `json:"type"`
	Label       string        `json:"label"`
	Description string        `json:"description"`
	Required    bool          `json:"required"`
	Default     interface{}   `json:"default,omitempty"`
	Options     []FieldOption `json:"options,omitempty"`
	Validation  *Validation   `json:"validation,omitempty"`
	Sensitive   bool          `json:"sensitive"`
	Placeholder string        `json:"placeholder,omitempty"`
	Fields      []ConfigField `json:"fields,omitempty"`
	IsArray     bool          `json:"is_array,omitempty"`
	ShowWhen    *ShowWhen     `json:"show_when,omitempty"`
}

type FieldOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Validation struct {
	Pattern string `json:"pattern,omitempty"`
	Min     *int   `json:"min,omitempty"`
	Max     *int   `json:"max,omitempty"`
	MinLen  *int   `json:"min_len,omitempty"`
	MaxLen  *int   `json:"max_len,omitempty"`
}

func GenerateConfigSpec(configType interface{}) []ConfigField {
	return generateConfigSpecRecursive(configType, "")
}

func generateConfigSpecRecursive(configType interface{}, prefix string) []ConfigField {
	var fields []ConfigField

	t := reflect.TypeOf(configType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		jsonTag := field.Tag.Get("json")

		// Handle inline embedded structs by recursively processing their fields
		if jsonTag != "" && strings.Contains(jsonTag, "inline") {
			embeddedType := field.Type
			if embeddedType.Kind() == reflect.Ptr {
				embeddedType = embeddedType.Elem()
			}

			// Create a zero value instance and recursively process it
			if embeddedType.Kind() == reflect.Struct {
				embeddedInstance := reflect.New(embeddedType).Interface()
				embeddedFields := generateConfigSpecRecursive(embeddedInstance, prefix)
				fields = append(fields, embeddedFields...)
			}
			continue
		}

		if field.Anonymous {
			continue
		}

		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		jsonName := strings.Split(jsonTag, ",")[0]

		description := field.Tag.Get("description")
		sensitive := field.Tag.Get("sensitive") == "true"
		defaultValue := field.Tag.Get("default")
		validateTag := field.Tag.Get("validate")
		placeholder := field.Tag.Get("placeholder")

		fieldType := inferFieldType(field.Type, sensitive)

		required := false
		if validateTag != "" {
			required = strings.Contains(validateTag, "required")
		}

		label := field.Tag.Get("label")
		if label == "" {
			label = toLabel(jsonName)
		}

		configField := ConfigField{
			Name:        jsonName,
			Type:        fieldType,
			Label:       label,
			Description: description,
			Required:    required,
			Sensitive:   sensitive,
			Placeholder: placeholder,
		}

		if defaultValue != "" {
			configField.Default = parseDefault(defaultValue, field.Type)
		}

		if validateTag != "" && strings.Contains(validateTag, "oneof=") {
			options := parseOneOfOptions(validateTag)
			if len(options) > 0 {
				configField.Type = FieldTypeSelect
				configField.Options = options
			}
		}

		showWhenTag := field.Tag.Get("show_when")
		if showWhenTag != "" {
			parts := strings.SplitN(showWhenTag, ":", 2)
			if len(parts) == 2 {
				configField.ShowWhen = &ShowWhen{
					Field: parts[0],
					Value: parts[1],
				}
			}
		}

		if fieldType == FieldTypeObject {
			nestedType := field.Type
			if nestedType.Kind() == reflect.Ptr {
				nestedType = nestedType.Elem()
			}

			if nestedType.Kind() == reflect.Slice {
				elemType := nestedType.Elem()
				if elemType.Kind() == reflect.Struct {
					configField.IsArray = true
					nestedInstance := reflect.New(elemType).Interface()
					configField.Fields = generateConfigSpecRecursive(nestedInstance, prefix+jsonName+".")
				}
			} else if nestedType.Kind() == reflect.Struct {
				nestedInstance := reflect.New(nestedType).Interface()
				configField.Fields = generateConfigSpecRecursive(nestedInstance, prefix+jsonName+".")
			}
		}

		fields = append(fields, configField)
	}

	return fields
}

func parseOneOfOptions(validateTag string) []FieldOption {
	parts := strings.Split(validateTag, "oneof=")
	if len(parts) < 2 {
		return nil
	}

	oneofPart := parts[1]
	if idx := strings.Index(oneofPart, ","); idx != -1 {
		oneofPart = oneofPart[:idx]
	}

	values := strings.Fields(oneofPart)
	if len(values) == 0 {
		return nil
	}

	options := make([]FieldOption, 0, len(values))
	for _, value := range values {
		options = append(options, FieldOption{
			Label: toLabel(value),
			Value: value,
		})
	}

	return options
}

func inferFieldType(t reflect.Type, sensitive bool) FieldType {
	if sensitive {
		return FieldTypePassword
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return FieldTypeString
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return FieldTypeInt
	case reflect.Bool:
		return FieldTypeBool
	case reflect.Slice:
		// Check if it's a slice of structs (array of objects)
		elemType := t.Elem()
		if elemType.Kind() == reflect.Struct {
			return FieldTypeObject // Will be handled as array of objects
		}
		return FieldTypeMultiselect
	case reflect.Struct:
		return FieldTypeObject
	default:
		return FieldTypeString
	}
}

func toLabel(fieldName string) string {
	parts := strings.Split(fieldName, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[0:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}

func parseDefault(value string, t reflect.Type) interface{} {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val, err := strconv.ParseInt(value, 10, 64); err == nil {
			return val
		}
	case reflect.Bool:
		if val, err := strconv.ParseBool(value); err == nil {
			return val
		}
	case reflect.String:
		return value
	}
	return value
}

// GetSensitiveFields returns a map of sensitive field names from ConfigSpec
func GetSensitiveFields(configSpec []ConfigField) map[string]bool {
	sensitive := make(map[string]bool)
	for _, field := range configSpec {
		if field.Sensitive {
			sensitive[field.Name] = true
		}
	}
	return sensitive
}

func ValidateConfigAgainstSpec(config map[string]interface{}, spec []ConfigField) error {
	return validateConfigRecursive(config, spec, "")
}

func validateConfigRecursive(config map[string]interface{}, spec []ConfigField, prefix string) error {
	for _, field := range spec {
		fieldPath := field.Name
		if prefix != "" {
			fieldPath = prefix + "." + field.Name
		}

		value, exists := config[field.Name]

		if field.Required {
			if !exists || isEmptyValue(value) {
				return fmt.Errorf("%s is required", fieldPath)
			}
		}

		if !exists {
			continue
		}

		if isEmptyValue(value) {
			continue
		}

		if err := validateFieldType(value, field, fieldPath); err != nil {
			return err
		}

		if field.Type == FieldTypeObject && len(field.Fields) > 0 {
			if field.IsArray {
				arr, ok := value.([]interface{})
				if !ok {
					return fmt.Errorf("%s must be an array", fieldPath)
				}
				for i, item := range arr {
					itemMap, ok := item.(map[string]interface{})
					if !ok {
						return fmt.Errorf("%s[%d] must be an object", fieldPath, i)
					}
					if err := validateConfigRecursive(itemMap, field.Fields, fmt.Sprintf("%s[%d]", fieldPath, i)); err != nil {
						return err
					}
				}
			} else {
				nestedMap, ok := value.(map[string]interface{})
				if !ok {
					return fmt.Errorf("%s must be an object", fieldPath)
				}
				if err := validateConfigRecursive(nestedMap, field.Fields, fieldPath); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func isEmptyValue(v interface{}) bool {
	if v == nil {
		return true
	}

	switch val := v.(type) {
	case string:
		return val == ""
	case []interface{}:
		return len(val) == 0
	case map[string]interface{}:
		return len(val) == 0
	default:
		return false
	}
}

func validateFieldType(value interface{}, field ConfigField, fieldPath string) error {
	switch field.Type {
	case FieldTypeString, FieldTypePassword:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("%s must be a string", fieldPath)
		}
	case FieldTypeInt:
		switch value.(type) {
		case int, int64, float64:
			// JSON numbers can come as float64
		default:
			return fmt.Errorf("%s must be a number", fieldPath)
		}
	case FieldTypeBool:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("%s must be a boolean", fieldPath)
		}
	case FieldTypeSelect:
		strVal, ok := value.(string)
		if !ok {
			return fmt.Errorf("%s must be a string", fieldPath)
		}
		if len(field.Options) > 0 {
			valid := false
			for _, opt := range field.Options {
				if opt.Value == strVal {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("%s has invalid value", fieldPath)
			}
		}
	case FieldTypeMultiselect:
		arr, ok := value.([]interface{})
		if !ok {
			return fmt.Errorf("%s must be an array", fieldPath)
		}
		for i, item := range arr {
			if _, ok := item.(string); !ok {
				return fmt.Errorf("%s[%d] must be a string", fieldPath, i)
			}
		}
	}
	return nil
}
