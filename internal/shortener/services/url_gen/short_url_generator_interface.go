// Package urlgen Генерация коротких ссылок
package urlgen

// ShortURLGeneratorInterface Генератор коротких версий URL
//
//go:generate mockgen -source=short_url_generator_interface.go -destination=mock/short_url_generator.go -package=mock ShortURLGeneratorInterface
type ShortURLGeneratorInterface interface {
	// GenerateShortURL Генерирует короткую версию URL, например, HnsSMA
	GenerateShortURL() (string, error)
}
