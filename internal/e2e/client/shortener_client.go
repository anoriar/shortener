package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	dtoRequestPkg "github.com/anoriar/shortener/internal/e2e/client/dto/request"
	dtoResponsePkg "github.com/anoriar/shortener/internal/e2e/client/dto/response"
	"io"
	"net/http"
)

type ShortenerClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewShortenerClient(httpClient *http.Client, baseURL string) *ShortenerClient {
	return &ShortenerClient{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (client *ShortenerClient) AddURL(url string) (*dtoResponsePkg.AddResponseDto, error) {
	request, err := http.NewRequest(http.MethodPost, client.baseURL, bytes.NewReader([]byte(url)))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "text/plain")
	response, err := client.httpClient.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return dtoResponsePkg.NewShortenerResponseDto(
		response.StatusCode,
		response.Header.Get("Content-Type"),
		string(body),
	), nil
}

func (client *ShortenerClient) GetURL(key string) (*dtoResponsePkg.GetResponseDto, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", client.baseURL, key), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "text/plain")
	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return dtoResponsePkg.NewGetResponseDto(
		response.StatusCode,
		response.Header.Get("Location"),
	), nil
}

func (client *ShortenerClient) AddURLv2(url string) (*dtoResponsePkg.AddResponseV2Dto, error) {
	requestDto := dtoRequestPkg.AddURLRequestDto{URL: url}
	requestJSON, err := json.Marshal(requestDto)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, client.baseURL+"/api/shorten", bytes.NewReader(requestJSON))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "text/plain")
	resp, err := client.httpClient.Do(request)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(string(body))
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		return nil, errors.New("not expected content type in response")
	}

	var addURLResponseDto dtoResponsePkg.AddURLResponseDTO
	err = json.Unmarshal(body, &addURLResponseDto)
	if err != nil {
		return nil, err
	}

	return dtoResponsePkg.NewAddResponseV2Dto(
		resp.StatusCode,
		resp.Header.Get("Content-Type"),
		addURLResponseDto,
	), nil
}
