package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gustavopnhro/STUGolang/database"
	"github.com/Gustavopnhro/STUGolang/structs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func main() {
	var port string = ":8080"

	database.Initialize()

	http.HandleFunc("/items", handleItems)

	http.ListenAndServe(port, nil)

	fmt.Println("Rodando API na porta", port)
}

func handleItems(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		searchProduct(writer, request)
	case http.MethodPost:
		createProduct(writer, request)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("Método não permitido"))
	}
}

func NewProduct(name string, price float64, stock int) *structs.Product {
	return &structs.Product{
		ID:    uuid.New(),
		Name:  name,
		Price: price,
		Stock: stock,
	}

}

func searchProduct(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/items" {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	productName := request.URL.Query().Get("name")

	if productName == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Nome do produto é necessário"))
		return
	}

	db := database.GetDB()

	statement, err := db.Prepare("SELECT id, name, price FROM Products WHERE name = ?")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(fmt.Sprintf("Erro ao preparar a declaração SQL: %v", err)))
		return
	}
	defer statement.Close()

	rows, err := statement.Query(productName)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(fmt.Sprintf("Erro ao executar a declaração SQL: %v", err)))
		return
	}
	defer rows.Close()

	var products []structs.Product

	for rows.Next() {
		var product structs.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(fmt.Sprintf("Erro ao escanear linha do resultado: %v", err)))
			return
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(fmt.Sprintf("Erro ao iterar através do resultado: %v", err)))
		return
	}

	productJSON, err := json.Marshal(products)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(fmt.Sprintf("Erro ao converter produtos para JSON: %v", err)))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(productJSON)
}

func createProduct(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("Not Allowed"))
		return
	}

	err := request.ParseForm()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Erro ao ler corpo da requisição"))
		return
	}

	productName := request.PostFormValue("name")
	productStockStr := request.PostFormValue("stock")
	productPriceStr := request.PostFormValue("price")

	if productName == "" || productStockStr == "" || productPriceStr == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Campos 'name', 'stock' e 'price' são obrigatórios"))
		return
	}

	stock, err := strconv.Atoi(productStockStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("O campo 'stock' deve ser um número válido"))
		return
	}

	price, err := strconv.ParseFloat(productPriceStr, 64)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("O campo 'price' deve ser um número válido"))
		return
	}

	newProduct := NewProduct(productName, price, stock)

	db := database.GetDB()

	statement, err := db.Prepare("INSERT INTO Products(id, name, price, stock) VALUES (?, ?, ?, ?)")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Erro ao preparar a declaração SQL"))
		return
	}
	defer statement.Close()

	_, err = statement.Exec(newProduct.ID, newProduct.Name, newProduct.Price, newProduct.Stock)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Erro ao executar a declaração SQL"))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte("Produto criado com sucesso"))
}
