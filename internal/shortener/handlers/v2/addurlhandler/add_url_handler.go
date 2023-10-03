package addurlhandler

import (
	"encoding/json"
	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/dto/response"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository"
	urlgen "github.com/anoriar/shortener/internal/shortener/services/url_gen"
	"github.com/google/uuid"
	"io"
	"net/http"
	neturl "net/url"
)

type AddHandler struct {
	urlRepository     repository.URLRepositoryInterface
	shortURLGenerator urlgen.ShortURLGeneratorInterface
	baseURL           string
}

func NewAddHandler(urlRepository repository.URLRepositoryInterface, shortURLGenerator urlgen.ShortURLGeneratorInterface, baseURL string) *AddHandler {
	return &AddHandler{
		urlRepository:     urlRepository,
		shortURLGenerator: shortURLGenerator,
		baseURL:           baseURL,
	}
}

func Initialize(cnf *config.Config, urlRepository repository.URLRepositoryInterface) *AddHandler {
	return NewAddHandler(urlRepository, urlgen.InitializeShortURLGenerator(urlRepository), cnf.BaseURL)
}

func (handler AddHandler) AddURL(w http.ResponseWriter, req *http.Request) {

	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	addURLRequestDto := request.AddURLRequestDto{}

	if err = json.Unmarshal(requestBody, &addURLRequestDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parsedURL, err := neturl.Parse(addURLRequestDto.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		http.Error(w, "Not valid URL", http.StatusBadRequest)
		return
	}

	shortKey, err := handler.shortURLGenerator.GenerateShortURL()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = handler.urlRepository.AddURL(&entity.URL{
		UUID:        uuid.NewString(),
		ShortURL:    shortKey,
		OriginalURL: addURLRequestDto.URL,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseDTO := response.AddURLResponseDto{
		Result: handler.baseURL + "/" + shortKey,
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	jsonResult, err := json.Marshal(responseDTO)
	if err != nil {
		return
	}

	_, err = w.Write(jsonResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
