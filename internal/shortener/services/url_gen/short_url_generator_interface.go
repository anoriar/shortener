package urlgen

type ShortURLGeneratorInterface interface {
	GenerateShortURL() (string, error)
}
