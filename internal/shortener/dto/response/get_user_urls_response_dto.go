package response

// GetUserURLsResponseDTO Результат получения всех URL, созданных пользователем
type GetUserURLsResponseDTO struct {
	ShortURL    string `json:"short_url"`    // Готовая версия для перехода http://localhost:8080/s1uHaW
	OriginalURL string `json:"original_url"` // оригинальный URL https://practicum1.yandex.ru
}
