package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	dtoRequestPkg "github.com/anoriar/shortener/internal/e2e/client/dto/request"
	dtoResponsePkg "github.com/anoriar/shortener/internal/e2e/client/dto/response"
	"github.com/anoriar/shortener/internal/e2e/config"
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
func InitializeShortenerClient(cnf *config.TestConfig) ShortenerClientInterface {
	return &ShortenerClient{
		httpClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Disable automatic redirects
				return http.ErrUseLastResponse
			},
		},
		baseURL: cnf.BaseURL,
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

	request.Header.Add("Content-Type", "application/json")
	resp, err := client.httpClient.Do(request)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := resp.Body

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(string(body))
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		return nil, errors.New("not expected content type in responsewriter")
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

func (client *ShortenerClient) AddURLv2WithCompress(url string, contentType string) (*dtoResponsePkg.AddResponseV2EncodingDto, error) {
	requestDto := dtoRequestPkg.AddURLRequestDto{URL: url}
	requestJSON, err := json.Marshal(requestDto)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, client.baseURL+"/api/shorten", bytes.NewReader(requestJSON))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", contentType)
	request.Header.Add("Accept-Encoding", "gzip")
	resp, err := client.httpClient.Do(request)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer gzReader.Close()
		reader = gzReader
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(string(body))
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		return nil, errors.New("not expected content type in responsewriter")
	}

	var addURLResponseDto dtoResponsePkg.AddURLResponseDTO
	err = json.Unmarshal(body, &addURLResponseDto)
	if err != nil {
		return nil, err
	}

	return dtoResponsePkg.NewAddResponseV2EncodingDto(
		resp.StatusCode,
		resp.Header.Get("Content-Type"),
		resp.Header.Get("Content-Encoding"),
		addURLResponseDto,
	), nil
}
