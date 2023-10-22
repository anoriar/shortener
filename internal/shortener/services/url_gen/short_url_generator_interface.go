package urlgen

//go:generate mockgen -source=short_url_generator_interface.go -destination=mock/short_url_generator.go -package=mock ShortURLGeneratorInterface
type ShortURLGeneratorInterface interface {
	GenerateShortURL() (string, error)
}
