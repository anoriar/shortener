package request

// AddURLRequestDto Запрос на сохранение оригинального урла
type AddURLRequestDto struct {
	URL string `json:"url"` // оригинальный URL
}
