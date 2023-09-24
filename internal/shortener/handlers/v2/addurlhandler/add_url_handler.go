package addurlhandler

import (
	"encoding/json"
	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/dto/request"
	"github.com/anoriar/shortener/internal/shortener/dto/response"
	"github.com/anoriar/shortener/internal/shortener/storage"
	"github.com/anoriar/shortener/internal/shortener/util"
	"github.com/asaskevich/govalidator"
	"io"
	"net/http"
)

type AddHandler struct {
	urlRepository storage.URLStorageInterface
	keyGen        util.KeyGenInterface
	baseURL       string
}

func (handler AddHandler) AddURL(w http.ResponseWriter, req *http.Request) {

	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	addUrlRequestDto := request.AddURLRequestDto{}

	if err = json.Unmarshal(requestBody, &addUrlRequestDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err = govalidator.ValidateStruct(&addUrlRequestDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortKey := handler.keyGen.Generate()
	err = handler.urlRepository.AddURL(addUrlRequestDto.Url, shortKey)

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

func NewAddHandler(urlRepository storage.URLStorageInterface, keyGen util.KeyGenInterface, baseURL string) *AddHandler {
	return &AddHandler{
		urlRepository: urlRepository,
		keyGen:        keyGen,
		baseURL:       baseURL,
	}
}

func Initialize(cnf *config.Config, storage storage.URLStorageInterface) *AddHandler {
	return NewAddHandler(storage, util.NewKeyGen(), cnf.BaseURL)
}
