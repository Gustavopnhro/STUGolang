package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gustavopnhro/STUGolang/database"
	"github.com/Gustavopnhro/STUGolang/structs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func NewProduct(name string, price float64, stock int) *structs.Product {
	return &structs.Product{
		ID:    uuid.New(),
		Name:  name,
		Price: price,
		Stock: stock,
	}

}

func main() {
	var port string = ":8080"

	database.Initialize()

	http.HandleFunc("/buscar_produto", searchProduct)
	http.HandleFunc("/items", createProduct)
	http.ListenAndServe(port, nil)

	fmt.Println("Rodando API na porta", port)

}

func searchProduct(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/buscar_produto" {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	product_parameter := request.URL.Query().Get("id")

	if product_parameter == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Olá mundo"))
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

	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte("Produto criado com sucesso"))
}
