package routes

import (
    "database/sql"
    "encoding/json"
    "fmt"
}

func main() {
	http.HandleFunc("/buscar_produto", searchProduct)
	http.ListenAndServe(":8080", nil)
}