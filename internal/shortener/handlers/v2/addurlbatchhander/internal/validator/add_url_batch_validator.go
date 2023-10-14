package validator

import (
	"fmt"
	"github.com/anoriar/shortener/internal/shortener/dto/request"
	neturl "net/url"
)

type AddURLBatchValidator struct {
}

func NewAddURLBatchValidator() *AddURLBatchValidator {
	return &AddURLBatchValidator{}
}

func (validator *AddURLBatchValidator) Validate(urls []request.AddURLBatchRequestDTO) error {
	if len(urls) == 0 {
		return fmt.Errorf("array of URL must be not empty")
	}
	for _, url := range urls {
		if url.CorrelationID == "" {
			return fmt.Errorf("correlation_id can not be empty")
		}

		parsedURL, err := neturl.Parse(url.OriginalURL)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			return fmt.Errorf("correlation_id = %s: not valid URL", url.CorrelationID)
		}
	}

	duplicates := validator.findDuplicateCorrelationIDs(urls)
	if len(duplicates) != 0 {
		return fmt.Errorf("duplicated correlation_id found %s", duplicates)
	}

	return nil
}

func (validator *AddURLBatchValidator) findDuplicateCorrelationIDs(urls []request.AddURLBatchRequestDTO) []string {
	uniqIDs := make(map[string]string, len(urls))
	var duplicates []string

	for _, url := range urls {
		if _, exists := uniqIDs[url.CorrelationID]; exists {
			duplicates = append(duplicates, url.CorrelationID)
			continue
		}
		uniqIDs[url.CorrelationID] = url.OriginalURL
	}
	return duplicates
}
