package gormbulk

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

//BulkUpdate update multiple records at once
// [objects]        Must be a slice of struct
// [excludeColumns] Columns you want to exclude from insert. You can omit if there is no column you want to exclude.
func BulkUpdate(db *gorm.DB, objects []interface{}, excludeColumns ...string) error {
	// Split records with specified size not to exceed Database parameter limit
	for _, objSet := range splitObjects(objects, 1) {
		if err := updateObjSet(db, objSet, excludeColumns...); err != nil {
			return err
		}
	}
	return nil
}

func updateObjSet(db *gorm.DB, objects []interface{}, excludeColumns ...string) error {
	if len(objects) == 0 {
		return nil
	}

	firstAttrs, err := ExtractMapValue(objects[0], excludeColumns)
	if err != nil {
		return err
	}

	attrSize := len(firstAttrs)

	// Scope to eventually run SQL
	mainScope := db.NewScope(objects[0])
	// Store placeholders for embedding variables
	placeholders := make([]string, 0, attrSize)
	whereCondition := make([]string, 0, attrSize)

	// Replace with database column name
	dbColumns := make([]string, 0, attrSize)
	for _, key := range sortedKeys(firstAttrs) {
		dbColumns = append(dbColumns, mainScope.Quote(key))
	}

	for _, obj := range objects {
		objAttrs, err := ExtractMapValue(obj, excludeColumns)
		if err != nil {
			return err
		}

		pks, err := ExtractPrymaryKeys(obj)
		if err != nil {
			return err
		}

		// If object sizes are different, SQL statement loses consistency
		if len(objAttrs) != attrSize {
			return errors.New("attribute sizes are inconsistent")
		}

		scope := db.NewScope(obj)
		// Append variables
		variables := make([]string, 0, attrSize)
		for _, key := range sortedKeys(objAttrs) {
			scope.AddToVars(objAttrs[key])
			variables = append(variables, key+" = ?")
		}

		valueQuery := strings.Join(variables, ", ")
		placeholders = append(placeholders, valueQuery)

		pknew := make([]string, 0, attrSize)
		for _, key := range sortedKeys(pks) {
			scope.AddToVars(pks[key])
			pknew = append(pknew, key+" = ?")
		}

		whereQuery := strings.Join(pknew, " AND ")
		whereCondition = append(whereCondition, whereQuery)

		// Also append variables to mainScope
		mainScope.SQLVars = append(mainScope.SQLVars, scope.SQLVars...)

	}

	updateOption := ""
	if val, ok := db.Get("gorm:update_option"); ok {
		strVal, ok := val.(string)
		if !ok {
			return errors.New("gorm:update_option should be a string")
		}
		updateOption = strVal
	}

	mainScope.Raw(fmt.Sprintf("UPDATE %s SET %s where %s %s",
		mainScope.QuotedTableName(),
		strings.Join(placeholders, ", "),
		strings.Join(whereCondition, ", "),
		updateOption,
	))

	return db.Exec(mainScope.SQL, mainScope.SQLVars...).Error
}
