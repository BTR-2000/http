package todo

import "fmt"

type JSONError struct {
	Error string `json:"error"`
}

func NewJsonError(text string, er error) *JSONError {
	str := fmt.Sprintf("Ошибка: %s - %v", text, er)
	return &JSONError{
		Error: str,
	}
}
