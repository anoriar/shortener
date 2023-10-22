package e2e

import (
	"github.com/anoriar/shortener/internal/e2e/client"
	"github.com/anoriar/shortener/internal/e2e/client/dto/request"
	"github.com/anoriar/shortener/internal/e2e/config"
	"github.com/caarlos0/env/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
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
	addResponse, err := shortenerClient.AddURL(testURL)
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

	getResponse, err := shortenerClient.GetURL(key)
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
	addResponse, err := shortenerClient.AddURLv2(testURL)
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

	getResponse, err := shortenerClient.GetURL(key)
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
	addResponse, err := shortenerClient.AddURLv2WithCompress(testURL, "application/json")
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

	getResponse, err := shortenerClient.GetURL(key)
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
	addResponse, err := shortenerClient.AddURLBatch(batchRequestItems)
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
		getResponse, err := shortenerClient.GetURL(shortKey)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusTemporaryRedirect, getResponse.StatusCode)
		mapItem, existed := correlationIDSHortKeyMap[correlationID]
		assert.True(t, existed)
		assert.Equal(t, mapItem.OriginalURL, getResponse.Location)
	}
}
