// Package repository Запрос на получение URL из хранилища
package repository

// Query Запрос на получение URL из хранилища
type Query struct {
	ShortURLs    []string // фильтр по коротким версиям URL
	OriginalURLs []string // фильтр по оригинальным версиям URL
}
