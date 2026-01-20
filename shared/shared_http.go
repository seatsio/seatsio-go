package shared

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"strings"
	"time"
)

type AdditionalHeader func(headers *map[string]string)

func ApiClient(secretKey string, baseUrl string, additionalHeaders ...AdditionalHeader) *req.Client {
	client := req.C().SetBaseURL(baseUrl).
		SetCommonBasicAuth(secretKey, "").
		SetTimeout(10*time.Second).
		SetCommonRetryCount(5).
		SetCommonRetryBackoffInterval(400*time.Millisecond, 10*time.Second).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			return err == nil && resp.StatusCode == 429
		})
	headers := make(map[string]string)
	for _, opt := range additionalHeaders {
		opt(&headers)
	}
	for key, value := range headers {
		client.SetCommonHeader(key, value)
	}
	return client
}

func AssertOk[T interface{}](result *req.Response, err error, data *T) (*T, error) {
	err = AssertOkWithoutResult(result, err)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func AssertOkNoBody(result *req.Response, err error) error {
	err = AssertOkWithoutResult(result, err)
	if err != nil {
		return err
	}
	return nil
}

func AssertOkMap[T interface{}](result *req.Response, err error, data map[string]T) (map[string]T, error) {
	err = AssertOkWithoutResult(result, err)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func AssertOkArray[T interface{}](result *req.Response, err error, data *[]T) ([]T, error) {
	err = AssertOkWithoutResult(result, err)
	if err != nil {
		return nil, err
	}
	return *data, nil
}

func AssertOkWithoutResult(result *req.Response, err error) error {
	if err != nil {
		return err
	}
	if !result.IsSuccessState() {
		if strings.Contains(result.GetHeader("content-type"), "application/json") {
			errorResponse := &SeatsioErrorResponse{}
			err := json.Unmarshal(result.Bytes(), errorResponse)
			if err != nil {
				return err
			}
			return &SeatsioError{Message: errorResponse.Errors[0].Message, Code: errorResponse.Errors[0].Code}
		} else {
			return fmt.Errorf("server returned error %v. Body: %v", result.StatusCode, string(result.Bytes()))
		}
	}
	return nil
}

func WithAdditionalHeader(key string, value string) AdditionalHeader {
	return func(headers *map[string]string) {
		(*headers)[key] = value
	}
}
