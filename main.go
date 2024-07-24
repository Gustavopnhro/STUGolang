package main

import (
	"database/sql"

	"github.com/Gustavopnhro/STUGolang/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func NewProduct(name string, price float64, stock int) *models.Product {
	return &models.Product{
		ID:    uuid.New().String(),
		Name:  name,
		Price: price,
		Stock: stock,
	}

}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/mysql")

	if err != nil {
		panic(err)
	}
	alicio := NewProduct("Alicio", 12.66, 11)
	err = insertCustomer(db, alicio)
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func insertCustomer(db *sql.DB, product *models.Product) error {
	statement, err := db.Prepare("insert into Customers(id, name, email, phone) values (?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}

	defer statement.Close()
	println(product.ID)

	_, err = statement.Exec(product.ID, product.Name, product.Stock, product.Price)
	if err != nil {
		return err
	}

	return nil
}
