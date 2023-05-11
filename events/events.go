package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"strings"
	"time"
)

type Events struct {
	secretKey string
	baseUrl   string
}

type createEventRequest struct {
	ChartKey string `json:"chartKey"`
}

type SeatsioError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SeatsioErrorResponse struct {
	Errors []SeatsioError `json:"errors"`
}

func ApiClient(secretKey string, baseUrl string) *req.Client {
	return req.C().SetBaseURL(baseUrl).
		SetCommonBasicAuth(secretKey, "").
		SetCommonRetryCount(5).
		SetCommonRetryBackoffInterval(400*time.Millisecond, 10*time.Second).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			return err == nil && resp.StatusCode == 429
		})
}

func (events *Events) Create(chartKey string) (*Event, error) {
	var event Event
	client := ApiClient(events.secretKey, events.baseUrl)
	result, err := client.R().
		SetBody(&createEventRequest{
			chartKey,
		}).
		SetSuccessResult(&event).
		Post("/events")
	return AssertOk(result, err, &event)
}

func AssertOk[T interface{}](result *req.Response, err error, data *T) (*T, error) {
	if err != nil {
		return nil, err
	}
	if result.IsErrorState() {
		if strings.Contains(result.GetHeader("content-type"), "application/json") {
			errorResponse := &SeatsioErrorResponse{}
			err := json.Unmarshal(result.Bytes(), errorResponse)
			if err != nil {
				return nil, err
			}
			return nil, errors.New(errorResponse.Errors[0].Message)
		} else {
			return nil, fmt.Errorf("server returned error %v. Body: %v", result.StatusCode, string(result.Bytes()))
		}
	}
	// TODO: what about 'unknown' state?
	return data, nil
}

func NewEvents(secretKey string, baseUrl string) *Events {
	return &Events{secretKey, baseUrl}
}
