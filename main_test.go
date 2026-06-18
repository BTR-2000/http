package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"proj/handlers"
	"proj/todo"
	"testing"
)

func TestCreate(t *testing.T) {
	store := todo.NewBookStore()
	h := handlers.NewHandler(store)

	bookDTO := todo.BookDTO{
		Title:  "Геракл",
		Author: "Агата Кристи",
		Year:   1950,
	}

	body, err := json.Marshal(bookDTO)
	if err != nil {
		t.Fatalf("Ошибка маршалинга в JSON %v", err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))

	r.Header.Set("Content-Type", "application/json")

	h.HandleCreateBook(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("Ожидался статус 201, получили %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Ожидался Content-Type - application/json, получили: %s", w.Header().Get("Content-Type"))
	}

	var book todo.Book
	if err = json.NewDecoder(w.Body).Decode(&book); err != nil {
		t.Fatalf("Ошибка декодирования JSON: %v", err)
	}

	if book.Author != "Агата Кристи" {
		t.Errorf("Ожидался автор - Агата Кристи, получили : %s", book.Author)
	}

	if book.Title != "Геракл" {
		t.Errorf("Ожидалось название - Геракл, получили : %s", book.Title)
	}

	if book.Year != 1950 {
		t.Errorf("Ожидался год: 1950, получили: %d", bookDTO.Year)
	}

	if book.ID != 1 {
		t.Errorf("Ожидался id = 1, получили: %d", book.ID)
	}

	resultBook, err := store.FindByID(1)
	if err != nil {
		t.Fatalf("Книга не найдена: %v", err)
	}
	if resultBook.Title != "Геракл" {
		t.Errorf("Название книг не совпадает. Ожидали 'Геракл', получили: %s", resultBook.Title)
	}
}
