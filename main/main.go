package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	blk "github.com/juliansarmiento9310/gorm-bulk-insert"
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
	var data Product
	data.Code = "ABC"
	data.Price = 10000

	var data2 Product
	data2.Code = "ABC1"
	data2.Price = 1000

	var data3 Product
	data3.Code = "ABC2"
	data3.Price = 2000

	var data4 Product
	data4.Code = "ABC3"
	data4.Price = 3000

	var data5 Product
	data5.Code = "ABC4"
	data5.Price = 4000

	var data6 Product
	data6.Code = "ABC5"
	data6.Price = 5000

	var massiveData []interface{}
	massiveData = append(massiveData, data)
	massiveData = append(massiveData, data2)
	massiveData = append(massiveData, data3)
	massiveData = append(massiveData, data4)
	massiveData = append(massiveData, data5)
	massiveData = append(massiveData, data6)
	fmt.Printf("data: ", massiveData)
	err = blk.BulkUpdate(db, massiveData, 1)

	var products []Product
	db.Find(&products)
	fmt.Printf("data: ", products[0])

}
