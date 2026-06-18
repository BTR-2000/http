package todo

import "fmt"

type BookDTO struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

func (b *BookDTO) Validate() error {
	if b.Title == "" {
		return fmt.Errorf("название книги не может быть пустым")
	}

	if b.Author == "" {
		return fmt.Errorf("имя автора не может быть пустым")
	}

	if b.Year < 1000 || b.Year > 2026 {
		return fmt.Errorf("Год выпуска книги должен быть между 1000 и 2026")
	}
	return nil
}
