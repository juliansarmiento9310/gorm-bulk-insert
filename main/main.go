package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	blk "github.com/juliansarmiento9310/gorm-bulk-insert"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	fmt.Println("here")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var data Product
	data.Code = "ABC"
	data.Price = 10000

	var massiveData []interface{}
	massiveData = append(massiveData, data)
	err = blk.BulkInsert(db, massiveData, 1)

}
