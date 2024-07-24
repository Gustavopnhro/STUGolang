package main

import (
	"database/sql"
	"models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gustavopnhro/module-database/scenario-1/models"
)

func NewCustomer(name string, email string, phone string) *models.Product {
	return &models.Product{
		ID:    uuid.New().String(),
		Name:  name,
		Email: email,
		Phone: phone,
	}

}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/mysql")

	if err != nil {
		panic(err)
	}
	alicio := NewCustomer("Alicio", "Alicio@gmail.com", "77889922")
	err = insertCustomer(db, alicio)
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func insertCustomer(db *sql.DB, customer *models.Product) error {
	statement, err := db.Prepare("insert into Customers(id, name, email, phone) values (?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}

	defer statement.Close()
	println(customer.ID)

	_, err = statement.Exec(customer.ID, customer.Name, customer.Email, customer.Phone)
	if err != nil {
		return err
	}

	return nil
}
