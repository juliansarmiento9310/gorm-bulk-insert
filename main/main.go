package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	blk "github.com/juliansarmiento9310/gorm-bulk-insert/v2"
)

type Product struct {
	gorm.Model
	Code  string `gorm:"primary_key"`
	Price uint
}

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var products []Product

	var data Product
	data.Code = "ABC"
	data.Price = 250

	var data2 Product
	data2.Code = "L1212"
	data2.Price = 121

	var massiveData []interface{}
	massiveData = append(massiveData, data)
	massiveData = append(massiveData, data2)
	err = blk.BulkUpdate(db, massiveData)
	fmt.Printf("-> ", err)

	db.Find(&products)

	fmt.Println("data: ", products[0].Price)
	fmt.Println("data: ", products[1].Price)
}
