package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	dtoRequestPkg "github.com/anoriar/shortener/internal/e2e/client/dto/request"
	dtoResponsePkg "github.com/anoriar/shortener/internal/e2e/client/dto/response"
	"github.com/anoriar/shortener/internal/e2e/config"
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

func (client *ShortenerClient) getTokenFromResponse(resp *http.Response) string {
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			return cookie.Value
		}
	}
	return ""
}

func (client *ShortenerClient) AddURL(requestDto dtoRequestPkg.AddURLRequestDto) (*dtoResponsePkg.AddResponseDto, error) {
	request, err := http.NewRequest(http.MethodPost, client.baseURL, bytes.NewReader([]byte(requestDto.URL)))
	if err != nil {
		return nil, err
	}
	if requestDto.IsAuth {
		request.AddCookie(&http.Cookie{
			Name:  "token",
			Value: requestDto.Token,
		})
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
		client.getTokenFromResponse(response),
		string(body),
	), nil
}

func (client *ShortenerClient) GetURL(requestDto dtoRequestPkg.GetURLRequestDto) (*dtoResponsePkg.GetResponseDto, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", client.baseURL, requestDto.ShortKey), nil)
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
		response.Header.Get("Content-Type"),
		response.Header.Get("Location"),
		client.getTokenFromResponse(response),
	), nil
}

func (client *ShortenerClient) AddURLv2(requestDto dtoRequestPkg.AddURLRequestDto) (*dtoResponsePkg.AddResponseV2Dto, error) {
	requestJSON, err := json.Marshal(struct {
		URL string `json:"url"`
	}{
		URL: requestDto.URL,
	})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, client.baseURL+"/api/shorten", bytes.NewReader(requestJSON))
	if err != nil {
		return nil, err
	}
	if requestDto.IsAuth {
		request.AddCookie(&http.Cookie{
			Name:  "token",
			Value: requestDto.Token,
		})
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
		client.getTokenFromResponse(resp),
		addURLResponseDto,
	), nil
}

func (client *ShortenerClient) AddURLv2WithCompress(requestDto dtoRequestPkg.AddURLRequestDto, contentType string) (*dtoResponsePkg.AddResponseV2EncodingDto, error) {
	requestJSON, err := json.Marshal(struct {
		URL string `json:"url"`
	}{
		URL: requestDto.URL,
	})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, client.baseURL+"/api/shorten", bytes.NewReader(requestJSON))
	if err != nil {
		return nil, err
	}
	if requestDto.IsAuth {
		request.AddCookie(&http.Cookie{
			Name:  "token",
			Value: requestDto.Token,
		})
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
		client.getTokenFromResponse(resp),
		resp.Header.Get("Content-Encoding"),
		addURLResponseDto,
	), nil
}

func (client *ShortenerClient) Ping() (dtoResponsePkg.PingResponseDto, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/ping", client.baseURL), nil)
	if err != nil {
		return dtoResponsePkg.PingResponseDto{}, err
	}

	request.Header.Add("Content-Type", "text/plain")
	response, err := client.httpClient.Do(request)
	if err != nil {
		return dtoResponsePkg.PingResponseDto{}, err
	}
	defer response.Body.Close()

	return dtoResponsePkg.NewPingResponseDto(
		response.StatusCode,
	), nil
}

func (client *ShortenerClient) DeleteURLBatch(shortURLs []string) (*dtoResponsePkg.DeleteURLBatchResponseDto, error) {
	requestJSON, err := json.Marshal(shortURLs)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodDelete, client.baseURL+"/api/shorten/batch", bytes.NewReader(requestJSON))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	resp, err := client.httpClient.Do(request)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return dtoResponsePkg.NewDeleteURLBatchResponseDto(
		resp.StatusCode,
		resp.Header.Get("Content-Type"),
		client.getTokenFromResponse(resp),
	), nil
}

func (client *ShortenerClient) AddURLBatch(requestDto dtoRequestPkg.AddURLBatchRequestDTO) (*dtoResponsePkg.AddURLBatchResponseDto, error) {
	requestJSON, err := json.Marshal(requestDto.Items)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, client.baseURL+"/api/shorten/batch", bytes.NewReader(requestJSON))
	if err != nil {
		return nil, err
	}
	if requestDto.IsAuth {
		request.AddCookie(&http.Cookie{
			Name:  "token",
			Value: requestDto.Token,
		})
	}

	request.Header.Add("Content-Type", "application/json")
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
		return nil, errors.New("not expected content type in responsewriter")
	}

	var batchItems []dtoResponsePkg.AddURLBatchItemDTO
	err = json.Unmarshal(body, &batchItems)
	if err != nil {
		return nil, err
	}

	return dtoResponsePkg.NewAddURLBatchResponseDto(
		resp.StatusCode,
		resp.Header.Get("Content-Type"),
		client.getTokenFromResponse(resp),
		batchItems,
	), nil
}

func (client *ShortenerClient) GetUserURLs(requestDto dtoRequestPkg.AuthRequest) (*dtoResponsePkg.GetUserURLsResponseDto, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", client.baseURL, "api/user/urls"), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "text/plain")

	if requestDto.IsAuth {
		request.AddCookie(&http.Cookie{
			Name:  "token",
			Value: requestDto.Token,
		})
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.Header.Get("Content-Type") != "application/json" {
		return nil, errors.New("not expected content type in responsewriter")
	}

	var items []dtoResponsePkg.UserURLResponseItem
	err = json.Unmarshal(body, &items)
	if err != nil {
		return nil, err
	}

	return dtoResponsePkg.NewGetUserURLsResponseDto(
		response.StatusCode,
		response.Header.Get("Content-Type"),
		client.getTokenFromResponse(response),
		items,
	), nil
}

func (client *ShortenerClient) DeleteUserURLs(requestDto dtoRequestPkg.DeleteUserURLsRequestDto) (*dtoResponsePkg.DeleteUserURLsResponseDto, error) {
	requestJSON, err := json.Marshal(requestDto.ShortURLs)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", client.baseURL, "api/user/urls"), bytes.NewReader(requestJSON))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "text/plain")
	if requestDto.IsAuth {
		request.AddCookie(&http.Cookie{
			Name:  "token",
			Value: requestDto.Token,
		})
	}
	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return dtoResponsePkg.NewDeleteUserURLsResponseDto(
		response.StatusCode,
		response.Header.Get("Content-Type"),
		client.getTokenFromResponse(response),
	), nil
}
