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

	// проверяем количество элементов в BookStore
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/books", nil)
	h.HandleGetBooks(w, r)

	slice := []todo.Book{}
	err = json.NewDecoder(w.Body).Decode(&slice)
	if err != nil {
		t.Fatalf("Ошибка декодирования: %v", err)
	}
	if len(slice) != 1 {
		t.Fatalf("Количество книг ожидалось 1, получили: %d", len(slice))
	}
}

func TestWrongCreate(t *testing.T) {
	store := todo.NewBookStore()
	h := handlers.NewHandler(store)

	bookDTO := todo.BookDTO{
		Title:  "",
		Author: "Толстой",
		Year:   2000,
	}

	body, err := json.Marshal(bookDTO)
	if err != nil {
		t.Fatalf("Ошибка маршалинга: %v", err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")

	h.HandleCreateBook(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус код 400, получили: %d", w.Code)
	}
}

func TestGetAndPutBookByID(t *testing.T) {
	store := todo.NewBookStore()
	h := handlers.NewHandler(store)

	bookDTO := todo.BookDTO{
		Title:  "Стрелок",
		Author: "Стивен Кинг",
		Year:   2000,
	}

	body, err := json.Marshal(bookDTO)
	if err != nil {
		t.Fatalf("Ошибка маршалинга: %v", err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")

	h.HandleCreateBook(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("Ожидался код 201, получили: %d", w.Code)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/books/{id}", h.HandleGetBookByID)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/books/1", nil)

	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("Ожидался статус код 200, получили: %d", w.Code)
	}

	var book todo.Book
	if err = json.NewDecoder(w.Body).Decode(&book); err != nil {
		t.Fatalf("Ошибка декодирования данных: %v", err)
	}

	if book.Author != "Стивен Кинг" {
		t.Errorf("Автор книги должен быть 'Стивен Кинг', получили ответ: %s", book.Author)
	}
	if book.Title != "Стрелок" {
		t.Errorf("Название книги ожидалось '', получили: %s", book.Title)
	}
	if book.Year != 2000 {
		t.Errorf("Год книги ожидался 2000, получили: %d", book.Year)
	}

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/books/2", nil)

	mux.ServeHTTP(w, r)

	if w.Code == http.StatusOK {
		t.Errorf("Получили статус %d, ожидали 404", w.Code)
	}

	// меняем поле IsRead
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPut, "/books/1", nil)

	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("Ожидался статус код 200, получили: %d", w.Code)
	}

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/books/1", nil)
	mux.ServeHTTP(w, r)
	// получаем книгу
	book = todo.Book{}
	if err := json.NewDecoder(w.Body).Decode(&book); err != nil {
		t.Fatalf("Ошибка декодирования: %v", err)
	}
	// проверяем статус книги
	if book.IsRead != true {
		t.Errorf("Ожидался статус книги true, получили: false")
	}
}

func TestDeleteBookByID(t *testing.T) {
	store := todo.NewBookStore()
	h := handlers.NewHandler(store)

	bookDTO := todo.BookDTO{
		Title:  "Стрелок",
		Author: "Стивен Кинг",
		Year:   2000,
	}

	body, err := json.Marshal(bookDTO)
	if err != nil {
		t.Fatalf("Ошибка маршалинга: %v", err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")

	h.HandleCreateBook(w, r)
	if w.Code != http.StatusCreated {
		t.Fatalf("Ожидался статус код 201, получили: %d", w.Code)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/books/{id}", h.HandleGetBookByID)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodDelete, "/books/1", nil)

	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("Ожидался статус код 200, получили: %d", w.Code)
	}

	sliceBooks := []todo.Book{}
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/books", nil)

	h.HandleGetBooks(w, r)

	err = json.NewDecoder(w.Body).Decode(&sliceBooks)
	if len(sliceBooks) != 0 {
		t.Errorf("Ожидалось 0 книг в остатке , получили: %d", len(sliceBooks))
	}
}

func TestGetBooksByAuthorAndYear(t *testing.T) {
	store := todo.NewBookStore()
	h := handlers.NewHandler(store)

	bookDTO := todo.BookDTO{
		Title:  "Подросток",
		Author: "Достоевский",
		Year:   2000,
	}
	bookDTO2 := todo.BookDTO{
		Title:  "Бесы",
		Author: "Достоевский",
		Year:   2000,
	}

	sliceBookDTO := []todo.BookDTO{bookDTO, bookDTO2}
	for _, v := range sliceBookDTO {
		body, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("Ошибка маршалинга: %v", err)
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")

		h.HandleCreateBook(w, r)

		if w.Code != http.StatusCreated {
			t.Fatalf("Ожидался статус код 201, получили: %d", w.Code)
		}
	}
	// тестируем книги по годам
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/books/year/2000", nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /books/year/{year}", h.HandleBooksByYear) //////
	mux.HandleFunc("GET /books/author/{author}", h.HandleBooksByAuthor)

	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("Ожидался статус код 200, получили: %d", w.Code)
	}

	booksByYear := []todo.Book{}
	if err := json.NewDecoder(w.Body).Decode(&booksByYear); err != nil {
		t.Fatalf("Ошибка декодирования: %v", err)
	}

	if len(booksByYear) != 2 {
		t.Errorf("Ожидалось 2 книги , получили: %d", len(booksByYear))
	}

	// тестируем книги по автору
	booksByAuthor := []todo.Book{}
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/books/author/Достоевский", nil)

	mux.ServeHTTP(w, r)
	if err := json.NewDecoder(w.Body).Decode(&booksByAuthor); err != nil {
		t.Fatalf("Ошибка декодирования: %v", err)
	}

	if len(booksByAuthor) != 2 {
		t.Errorf("Ожидалось 2 книги , получили: %d", len(booksByAuthor))
	}
}
