package shared

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/require"
)

func failedResponse(ctx context.Context, method string) *req.Response {
	request := req.C().R().SetContext(ctx)
	request.Method = method
	return &req.Response{Request: request}
}

func TestRetryConditionDoesNotRetryWhenContextIsCancelled(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	connectionError := errors.New("connection reset")

	for _, method := range []string{http.MethodGet, http.MethodPost} {
		require.False(t, retryCondition(true)(failedResponse(ctx, method), connectionError), method)
	}
}

func TestRetryConditionRetriesConnectionErrorsWhileContextIsActive(t *testing.T) {
	t.Parallel()
	connectionError := errors.New("connection reset")

	for _, method := range []string{http.MethodGet, http.MethodPost} {
		require.True(t, retryCondition(true)(failedResponse(context.Background(), method), connectionError), method)
	}
}
