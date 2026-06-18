package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"proj/todo"
	"strconv"
)

type Handler struct {
	store *todo.BookStore
}

func NewHandler(newStore *todo.BookStore) *Handler {
	return &Handler{
		store: newStore,
	}
}

func (h *Handler) HandleGetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	books := h.store.GetAll()
	if err := json.NewEncoder(w).Encode(books); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(todo.NewJsonError("Ошибка кодирования JSON", err))
		return
	}
}

func (h *Handler) HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var bookDTO todo.BookDTO
	err := json.NewDecoder(r.Body).Decode(&bookDTO)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(todo.NewJsonError("Ошибка декодирования JSON", err))
		return
	}

	if err := bookDTO.Validate(); err != nil {
		http.Error(w, "Неверные данные!", http.StatusBadRequest)
		return
	}

	book := h.store.CreateBook(bookDTO)

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		fmt.Println("Не удалось отправить ответ клиенту:", err)
		return
	}
}

func (h *Handler) HandleGetBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(todo.NewJsonError("Ошибка id должно быть целым числом:", err))
		return
	}

	book, err := h.store.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(todo.NewJsonError("Ошибка: ", err))
		return
	}
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(book)
	case http.MethodPut:
		h.store.UpdateBookByID(id)
		json.NewEncoder(w).Encode(map[string]string{"Результат": "данные книги успешно изменены"})
	case http.MethodDelete:
		err = h.store.DeleteBookByID(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(todo.NewJsonError("Ошибка: ", err))
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"Результат": "книга успешно удалена"})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) HandleBooksByAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	author := r.PathValue("author")
	if author == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	books := h.store.GetBooksByAuthor(author)

	if err := json.NewEncoder(w).Encode(books); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(todo.NewJsonError("Ошибка: ", err))
		return
	}
}

func (h *Handler) HandleBooksByYear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	yearStr := r.PathValue("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(todo.NewJsonError("Ошибка - год должен быть целым числом: ", err))
		return
	}

	books := h.store.GetBooksByYear(year)
	if err := json.NewEncoder(w).Encode(books); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(todo.NewJsonError("Ошибка: ", err))
		return
	}
}
