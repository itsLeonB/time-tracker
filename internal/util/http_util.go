package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/rotisserie/eris"
)

func ConstructUrl(baseUrl string, targetUrl string, queryParams map[string]string) string {
	if len(queryParams) == 0 {
		return fmt.Sprintf("%s/%s", baseUrl, targetUrl)
	}

	params := url.Values{}

	for key, value := range queryParams {
		params.Set(key, value)
	}

	return fmt.Sprintf("%s/%s?%s", baseUrl, targetUrl, params.Encode())
}

type RequestObject struct {
	Client    *http.Client
	Ctx       context.Context
	Method    string
	TargetUrl string
	Body      io.Reader
	Headers   map[string]string
}

func MakeRequest[SuccessType any, ErrorType any](requestObj *RequestObject) (SuccessType, ErrorType, error) {
	var result SuccessType
	var errorResponse ErrorType

	request, err := http.NewRequestWithContext(requestObj.Ctx, requestObj.Method, requestObj.TargetUrl, requestObj.Body)
	if err != nil {
		return result, errorResponse, eris.Wrap(err, "failed to create request")
	}

	if requestObj.Headers != nil {
		for key, value := range requestObj.Headers {
			request.Header.Set(key, value)
		}
	}

	response, err := requestObj.Client.Do(request)
	if err != nil {
		return result, errorResponse, eris.Wrap(err, "failed to execute request")
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return result, errorResponse, eris.Wrap(err, "failed to read response body")
	}

	// Try decoding into the expected result type
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		var jsonErr *json.UnmarshalTypeError
		if errors.As(err, &jsonErr) {
			// Try decoding into the error type
			err2 := json.Unmarshal(bodyBytes, &errorResponse)
			if err2 != nil {
				return result, errorResponse, eris.Wrap(err2, "failed to decode failure response")
			}
			log.Printf("request failed with status code %d: %v\n", response.StatusCode, errorResponse)
			return result, errorResponse, nil
		}
		return result, errorResponse, eris.Wrap(err, "failed to decode response")
	}

	return result, errorResponse, nil
}
