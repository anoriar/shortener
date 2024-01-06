package e2e

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/anoriar/shortener/internal/e2e/client"
	"github.com/anoriar/shortener/internal/e2e/client/dto/request"
	response2 "github.com/anoriar/shortener/internal/e2e/client/dto/response"
	"github.com/anoriar/shortener/internal/e2e/config"
)

const testURL = "https://github.com/"

func Test_Shortener(t *testing.T) {
	cnf := config.NewTestConfig()
	err := env.Parse(cnf)
	assert.NoError(t, err)

	if cnf.BaseURL == "" {
		t.Skip()
	}

	shortenerClient := client.InitializeShortenerClient(cnf)
	addResponse, err := shortenerClient.AddURL(request.AddURLRequestDto{
		URL: testURL,
	})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, addResponse.StatusCode)
	assert.Equal(t, "text/plain", addResponse.ContentType)
	assert.True(t, strings.HasPrefix(addResponse.Body, cnf.BaseURL))

	splittedURL := strings.Split(addResponse.Body, "/")
	key := splittedURL[len(splittedURL)-1]
	defer func() {
		_, err = shortenerClient.DeleteURLBatch([]string{key})
		require.NoError(t, err)
	}()

	getResponse, err := shortenerClient.GetURL(request.GetURLRequestDto{
		ShortKey: key,
	})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusTemporaryRedirect, getResponse.StatusCode)
	assert.Equal(t, testURL, getResponse.Location)
}

func Test_ShortenerV2(t *testing.T) {
	cnf := config.NewTestConfig()
	err := env.Parse(cnf)
	assert.NoError(t, err)

	if cnf.BaseURL == "" {
		t.Skip()
	}

	shortenerClient := client.InitializeShortenerClient(cnf)
	addResponse, err := shortenerClient.AddURLv2(request.AddURLRequestDto{
		URL: testURL,
	})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, addResponse.StatusCode)
	assert.Equal(t, "application/json", addResponse.ContentType)
	assert.True(t, strings.HasPrefix(addResponse.Body.Result, cnf.BaseURL))

	splittedURL := strings.Split(addResponse.Body.Result, "/")
	key := splittedURL[len(splittedURL)-1]
	defer func() {
		_, err = shortenerClient.DeleteURLBatch([]string{key})
		require.NoError(t, err)
	}()

	getResponse, err := shortenerClient.GetURL(request.GetURLRequestDto{
		ShortKey: key,
	})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusTemporaryRedirect, getResponse.StatusCode)
	assert.Equal(t, testURL, getResponse.Location)
}

func Test_ShortenerV2WithCompress(t *testing.T) {
	cnf := config.NewTestConfig()
	err := env.Parse(cnf)
	assert.NoError(t, err)

	if cnf.BaseURL == "" {
		t.Skip()
	}

	shortenerClient := client.InitializeShortenerClient(cnf)
	addResponse, err := shortenerClient.AddURLv2WithCompress(request.AddURLRequestDto{
		URL: testURL,
	}, "application/json")
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, addResponse.StatusCode)
	assert.Equal(t, "application/json", addResponse.ContentType)
	assert.Equal(t, "gzip", addResponse.ContentEncoding)
	assert.True(t, strings.HasPrefix(addResponse.Body.Result, cnf.BaseURL))

	splittedURL := strings.Split(addResponse.Body.Result, "/")
	key := splittedURL[len(splittedURL)-1]
	defer func() {
		_, err = shortenerClient.DeleteURLBatch([]string{key})
		require.NoError(t, err)
	}()

	getResponse, err := shortenerClient.GetURL(request.GetURLRequestDto{
		ShortKey: key,
	})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusTemporaryRedirect, getResponse.StatusCode)
	assert.Equal(t, testURL, getResponse.Location)

}

func Test_ShortenerPing(t *testing.T) {
	cnf := config.NewTestConfig()
	err := env.Parse(cnf)
	assert.NoError(t, err)

	if cnf.BaseURL == "" {
		t.Skip()
	}

	shortenerClient := client.InitializeShortenerClient(cnf)
	pingResponse, err := shortenerClient.Ping()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, pingResponse.StatusCode)
}

const (
	originalURL1   = "https://practicum.yandex.ru"
	correlationID1 = "g0fsdf9fj"
	originalURL2   = "https://practicum2.yandex.ru"
	correlationID2 = "ngfdsf3"
	originalURL3   = "https://practicum3.yandex.ru"
	correlationID3 = "by4564trg"
)

func Test_ShortenerAddURlBatch(t *testing.T) {
	cnf := config.NewTestConfig()
	err := env.Parse(cnf)
	assert.NoError(t, err)

	if cnf.BaseURL == "" {
		t.Skip()
	}

	correlationIDSHortKeyMap := map[string]request.AddURLBatchItemDTO{
		correlationID1: {
			CorrelationID: correlationID1,
			OriginalURL:   originalURL1,
		},
		correlationID2: {
			CorrelationID: correlationID2,
			OriginalURL:   originalURL2,
		},
		correlationID3: {
			CorrelationID: correlationID3,
			OriginalURL:   originalURL3,
		},
	}
	batchRequestItems := make([]request.AddURLBatchItemDTO, len(correlationIDSHortKeyMap))
	i := 0
	for _, mapItem := range correlationIDSHortKeyMap {
		batchRequestItems[i] = mapItem
		i++
	}

	shortenerClient := client.InitializeShortenerClient(cnf)
	addResponse, err := shortenerClient.AddURLBatch(request.AddURLBatchRequestDTO{
		Items: batchRequestItems,
	})
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, addResponse.StatusCode)
	assert.Equal(t, "application/json", addResponse.ContentType)
	assert.True(t, len(addResponse.Body) > 0)

	correlationIDShortKeyMap := make(map[string]string)
	for _, item := range addResponse.Body {
		assert.True(t, strings.HasPrefix(item.ShortURL, cnf.BaseURL))
		splittedURL := strings.Split(item.ShortURL, "/")
		correlationIDShortKeyMap[item.CorrelationID] = splittedURL[len(splittedURL)-1]
	}
	defer func() {
		keysForDelete := make([]string, len(correlationIDShortKeyMap))
		i := 0
		for _, key := range correlationIDShortKeyMap {
			keysForDelete[i] = key
			i++
		}
		_, err = shortenerClient.DeleteURLBatch(keysForDelete)
		require.NoError(t, err)
	}()

	for correlationID, shortKey := range correlationIDShortKeyMap {
		getResponse, err := shortenerClient.GetURL(request.GetURLRequestDto{
			ShortKey: shortKey,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusTemporaryRedirect, getResponse.StatusCode)
		mapItem, existed := correlationIDSHortKeyMap[correlationID]
		assert.True(t, existed)
		assert.Equal(t, mapItem.OriginalURL, getResponse.Location)
	}
}

func Test_ShortenerGetUserURLs(t *testing.T) {
	cnf := config.NewTestConfig()
	err := env.Parse(cnf)
	assert.NoError(t, err)

	if cnf.BaseURL == "" {
		t.Skip()
	}

	shortenerClient := client.InitializeShortenerClient(cnf)

	var expectedURLs []response2.UserURLResponseItem
	var keysForDelete []string
	originalURLs := []string{originalURL1, originalURL2, originalURL3}

	auth := &request.AuthRequest{
		Token:  "",
		IsAuth: false,
	}
	for i, url := range originalURLs {
		addResponse, err := shortenerClient.AddURLv2(request.AddURLRequestDto{
			AuthRequest: *auth,
			URL:         url,
		})
		assert.NoError(t, err)

		if i == 0 {
			auth.IsAuth = true
			auth.Token = addResponse.Token

		}

		splittedURL := strings.Split(addResponse.Body.Result, "/")
		keysForDelete = append(keysForDelete, splittedURL[len(splittedURL)-1])
		expectedURLs = append(expectedURLs, response2.UserURLResponseItem{
			ShortURL:    addResponse.Body.Result,
			OriginalURL: url,
		})
	}

	defer func() {
		_, err = shortenerClient.DeleteURLBatch(keysForDelete)
		require.NoError(t, err)
	}()

	response, err := shortenerClient.GetUserURLs(*auth)
	require.NoError(t, err)

	assert.True(t, len(expectedURLs) == len(response.Items))
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, expectedURLs, response.Items)
}

func Test_ShortenerDeleteUserURLs(t *testing.T) {
	cnf := config.NewTestConfig()
	err := env.Parse(cnf)
	assert.NoError(t, err)

	if cnf.BaseURL == "" {
		t.Skip()
	}

	shortenerClient := client.InitializeShortenerClient(cnf)

	var keysForDelete []string
	originalURLs := []string{originalURL1, originalURL2, originalURL3}

	auth := &request.AuthRequest{
		Token:  "",
		IsAuth: false,
	}
	for i, url := range originalURLs {
		addResponse, err := shortenerClient.AddURLv2(request.AddURLRequestDto{
			AuthRequest: *auth,
			URL:         url,
		})
		assert.NoError(t, err)

		if i == 0 {
			auth.IsAuth = true
			auth.Token = addResponse.Token
		}

		splittedURL := strings.Split(addResponse.Body.Result, "/")
		keysForDelete = append(keysForDelete, splittedURL[len(splittedURL)-1])
	}

	defer func() {
		_, err = shortenerClient.DeleteURLBatch(keysForDelete)
		require.NoError(t, err)
	}()

	response, err := shortenerClient.DeleteUserURLs(request.DeleteUserURLsRequestDto{
		AuthRequest: *auth,
		ShortURLs:   keysForDelete,
	})
	require.NoError(t, err)

	//Операция асинхронная
	time.Sleep(1 * time.Second)

	assert.Equal(t, http.StatusAccepted, response.StatusCode)

	for _, key := range keysForDelete {
		response, err := shortenerClient.GetURL(request.GetURLRequestDto{
			ShortKey: key,
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusGone, response.StatusCode)
	}
}
