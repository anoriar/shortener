package addurlhandler

import (
	"encoding/json"
	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/dto/response"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository"
	"github.com/anoriar/shortener/internal/shortener/util"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"io"
	"net/http"
)

type AddHandler struct {
	urlRepository repository.URLRepositoryInterface
	keyGen        util.KeyGenInterface
	baseURL       string
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

	if _, err = govalidator.ValidateStruct(&addURLRequestDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortKey := handler.keyGen.Generate()
	//TODO: доделать ы keygenerator генерацию в случае существюущего ключа (5 попыток)

	_, err = handler.urlRepository.AddURL(&entity.Url{
		Uuid:        uuid.NewString(),
		ShortURL:    shortKey,
		OriginalURL: addURLRequestDto.URL,
	})

	if err != nil {
		//TODO: middleware на обработку ошибок
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

func NewAddHandler(urlRepository repository.URLRepositoryInterface, keyGen util.KeyGenInterface, baseURL string) *AddHandler {
	return &AddHandler{
		urlRepository: urlRepository,
		keyGen:        keyGen,
		baseURL:       baseURL,
	}
}

func Initialize(cnf *config.Config, urlRepository repository.URLRepositoryInterface) *AddHandler {
	return NewAddHandler(urlRepository, util.NewKeyGen(), cnf.BaseURL)
}
