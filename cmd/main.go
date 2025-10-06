// main.go
package main

import (
	"fmt"
	"net/http"

	_ "task/docs"
	"task/pkg/server"
)

func main() {
	app := server.App()

	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	fmt.Println("Server is listening on port 8081")
	fmt.Println("Swagger UI: http://localhost:8081/swagger/index.html")
	fmt.Println("Swagger JSON: http://localhost:8081/swagger/doc.json")

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
