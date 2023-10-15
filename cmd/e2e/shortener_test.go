package e2e

import (
	"github.com/anoriar/shortener/internal/e2e/client"
	"github.com/anoriar/shortener/internal/e2e/config"
	"github.com/caarlos0/env/v6"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

const testURL = "https://github.com/"

// #MENTOR как дальше писать e2e тест? Нужно вседа пересоздавать базу данных, иначе возникают частые конфликты 409.
// 1 раз запустишь и каждый раз руками дропать бд.
// А возможно в локальной разработке она пригодится
// В идеале - создавать тестовую бдху, но тогда надо запускать и свой тестовые сервер отдельно
// через httptest.NewServer -? тогда надо переиницилазировать весь хендлер и по сути выполнять всю внутрянку main.go сервера + зависимости
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
