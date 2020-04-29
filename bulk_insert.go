package gormbulk

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

// Insert multiple records at once
// [objects]        Must be a slice of struct
// [chunkSize]      Number of records to insert at once.
//                  Embedding a large number of variables at once will raise an error beyond the limit of prepared statement.
//                  Larger size will normally lead the better performance, but 2000 to 3000 is reasonable.
// [excludeColumns] Columns you want to exclude from insert. You can omit if there is no column you want to exclude.
func BulkInsert(db *gorm.DB, objects []interface{}, chunkSize int, excludeColumns ...string) error {
	// Split records with specified size not to exceed Database parameter limit
	for _, objSet := range splitObjects(objects, chunkSize) {
		value, err := insertObjSet(db, objSet, excludeColumns...)
		if err != nil {
			return err
		}
		fmt.Printf("Valor retornado %v", value)
	}
	return nil
}

func insertObjSet(db *gorm.DB, objects []interface{}, excludeColumns ...string) (interface{}, error) {
	if len(objects) == 0 {
		return nil, nil
	}

	firstAttrs, err := ExtractMapValue(objects[0], excludeColumns)
	if err != nil {
		return nil, err
	}

	attrSize := len(firstAttrs)

	// Scope to eventually run SQL
	mainScope := db.NewScope(objects[0])
	// Store placeholders for embedding variables
	placeholders := make([]string, 0, attrSize)

	// Replace with database column name
	dbColumns := make([]string, 0, attrSize)
	for _, key := range sortedKeys(firstAttrs) {
		dbColumns = append(dbColumns, mainScope.Quote(key))
	}

	for _, obj := range objects {
		objAttrs, err := ExtractMapValue(obj, excludeColumns)
		if err != nil {
			return nil, err
		}

		// If object sizes are different, SQL statement loses consistency
		if len(objAttrs) != attrSize {
			return nil, errors.New("attribute sizes are inconsistent")
		}

		scope := db.NewScope(obj)

		// Append variables
		variables := make([]string, 0, attrSize)
		for _, key := range sortedKeys(objAttrs) {
			scope.AddToVars(objAttrs[key])
			variables = append(variables, "?")
		}

		valueQuery := "(" + strings.Join(variables, ", ") + ")"
		placeholders = append(placeholders, valueQuery)

		// Also append variables to mainScope
		mainScope.SQLVars = append(mainScope.SQLVars, scope.SQLVars...)
	}

	insertOption := ""
	if val, ok := db.Get("gorm:insert_option"); ok {
		strVal, ok := val.(string)
		if !ok {
			return nil, errors.New("gorm:insert_option should be a string")
		}
		insertOption = strVal
	}

	mainScope.Raw(fmt.Sprintf("INSERT INTO %s (%s) VALUES %s %s",
		mainScope.QuotedTableName(),
		strings.Join(dbColumns, ", "),
		strings.Join(placeholders, ", "),
		insertOption,
	))

	result := db.Exec(mainScope.SQL, mainScope.SQLVars...)

	return result.Value, result.Error
}
