package gormbulk

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

func ExtractPrymaryKeys(value interface{}) (map[string]interface{}, error) {
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		value = rv.Interface()
	}
	if rv.Kind() != reflect.Struct {
		return nil, errors.New("value must be kind of Struct")
	}

	var pksFields = map[string]interface{}{}
	for _, field := range (&gorm.Scope{Value: value}).Fields() {
		// Exclude relational record because it's not directly contained in database columns
		if fieldIsValidPrimary(field) {
			pksFields[field.DBName] = field.Field.Interface()
		}
	}
	return pksFields, nil
}

// Obtain columns and values required for insert from interface
func ExtractMapValue(value interface{}, excludeColumns []string) (map[string]interface{}, error) {
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		value = rv.Interface()
	}
	if rv.Kind() != reflect.Struct {
		return nil, errors.New("value must be kind of Struct")
	}

	var attrs = map[string]interface{}{}

	for _, field := range (&gorm.Scope{Value: value}).Fields() {
		// Exclude relational record because it's not directly contained in database columns
		_, hasForeignKey := field.TagSettingsGet("FOREIGNKEY")

		if !containString(excludeColumns, field.Struct.Name) && field.StructField.Relationship == nil && !hasForeignKey &&
			!field.IsIgnored && !fieldIsAutoIncrement(field) && !fieldIsPrimaryAndBlank(field) {
			if field.Struct.Name == "UpdatedAt" && field.IsBlank {
				attrs[field.DBName] = time.Now()
			} else if !fieldIsValidPrimary(field) && !field.IsBlank {
				attrs[field.DBName] = field.Field.Interface()
			}
		}
	}
	return attrs, nil
}

func fieldIsAutoIncrement(field *gorm.Field) bool {
	if value, ok := field.TagSettingsGet("AUTO_INCREMENT"); ok {
		return strings.ToLower(value) != "false"
	}
	return false
}

func fieldIsPrimaryAndBlank(field *gorm.Field) bool {
	return field.IsPrimaryKey && field.IsBlank
}

func fieldIsValidPrimary(field *gorm.Field) bool {
	return field.IsPrimaryKey && !field.IsBlank
}
