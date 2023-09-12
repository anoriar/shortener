package client

import (
	"bytes"
	"fmt"
	"github.com/anoriar/shortener/cmd/e2e/client/dto"
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

func (client *ShortenerClient) AddURL(url string) (*dto.AddResponseDto, error) {
	request, err := http.NewRequest(http.MethodPost, client.baseURL, bytes.NewReader([]byte(url)))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "text/plain")
	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return dto.NewShortenerResponseDto(
		response.StatusCode,
		response.Header.Get("Content-Type"),
		string(body),
	), nil
}

func (client *ShortenerClient) GetURL(key string) (*dto.GetResponseDto, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", client.baseURL, key), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "text/plain")
	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return dto.NewGetResponseDto(
		response.StatusCode,
		response.Header.Get("Location"),
	), nil
}
