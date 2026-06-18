package todo

import (
	"fmt"
	"sync"
	"time"
)

type Book struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Year      int    `json:"year"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

type BookStore struct {
	mu     sync.RWMutex
	books  map[int]Book
	nextID int
}

func NewBookStore() *BookStore {
	return &BookStore{
		books:  make(map[int]Book),
		nextID: 1,
	}
}

func (b *BookStore) GetAll() []Book {
	b.mu.RLock()
	defer b.mu.RUnlock()
	slice := make([]Book, 0, len(b.books))
	for _, book := range b.books {
		slice = append(slice, book)
	}
	return slice
}

func (b *BookStore) CreateBook(bookDTO BookDTO) Book {
	b.mu.Lock()
	defer b.mu.Unlock()
	book := Book{
		ID:        b.nextID,
		Title:     bookDTO.Title,
		Author:    bookDTO.Author,
		Year:      bookDTO.Year,
		CreatedAt: fmt.Sprintf("%s", time.Now().Format("2006-01-02 15:04:05")),
	}

	b.books[b.nextID] = book
	b.nextID++
	return book
}

func (b *BookStore) FindByID(id int) (Book, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	book, exist := b.books[id]
	if !exist {
		return Book{}, fmt.Errorf("Не существует книги с id = %d", id)
	}
	return book, nil
}

func (b *BookStore) UpdateBookByID(id int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	book := b.books[id]
	if book.IsRead {
		book.IsRead = false
	} else {
		book.IsRead = true
	}
	b.books[id] = book
}

func (b *BookStore) DeleteBookByID(id int) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, exist := b.books[id]
	if !exist {
		return fmt.Errorf("Не существует книги с id = %d", id)
	}
	delete(b.books, id)
	return nil
}

func (b *BookStore) GetBooksByAuthor(author string) []Book {
	b.mu.RLock()
	defer b.mu.RUnlock()
	books := make([]Book, 0)
	for _, book := range b.books {
		if book.Author == author {
			books = append(books, book)
		}
	}
	return books
}

func (b *BookStore) GetBooksByYear(year int) []Book {
	b.mu.RLock()
	defer b.mu.RUnlock()
	books := make([]Book, 0)
	for _, book := range b.books {
		if book.Year == year {
			books = append(books, book)
		}
	}
	return books
}
