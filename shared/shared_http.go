package shared

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"strings"
	"time"
)

type AdditionalConfig func(client *req.Client)

func ApiClient(secretKey string, baseUrl string, additionalConfig ...AdditionalConfig) *req.Client {
	client := req.C().SetBaseURL(baseUrl).
		SetCommonBasicAuth(secretKey, "").
		SetCommonRetryCount(5).
		SetCommonRetryBackoffInterval(400*time.Millisecond, 10*time.Second).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			return err == nil && resp.StatusCode == 429
		})
	for _, opt := range additionalConfig {
		opt(client)
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

func AdditionalHeader(key string, value string) AdditionalConfig {
	return func(client *req.Client) {
		client.SetCommonHeader(key, value)
	}
}
