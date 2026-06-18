package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"proj/handlers"
	"proj/todo"
)

func main() {
	store := todo.NewBookStore()
	h := handlers.NewHandler(store)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /books", h.HandleGetBooks)
	mux.HandleFunc("POST /books", h.HandleCreateBook)
	mux.HandleFunc("/books/{id}", h.HandleGetBookByID)
	mux.HandleFunc("/books/author/{author}", h.HandleBooksByAuthor)
	mux.HandleFunc("/books/year/{year}", h.HandleBooksByYear)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
	fmt.Println("Завершение программы")
}
