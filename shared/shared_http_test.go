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

	require.False(t, retryCondition(failedResponse(ctx, http.MethodGet), errors.New("connection reset")))
}

func TestRetryConditionRetriesGetsOnConnectionErrorsWhileContextIsActive(t *testing.T) {
	t.Parallel()

	require.True(t, retryCondition(failedResponse(context.Background(), http.MethodGet), errors.New("connection reset")))
}

func TestRetryConditionNeverRetriesWrites(t *testing.T) {
	t.Parallel()

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		require.False(t, retryCondition(failedResponse(context.Background(), method), errors.New("connection reset")), method)
	}
}

func TestIdempotencyKeyStatusOf(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		status  idempotencyKeyStatus
		known   bool
		honored bool
	}{
		"executed":    {statusExecuted, true, true},
		"replayed":    {statusReplayed, true, true},
		"in-use":      {statusInUse, true, true},
		"reused":      {statusReused, true, true},
		"ignored":     {statusIgnored, true, false},
		" IGNORED ":   {statusIgnored, true, false},
		"Executed":    {statusExecuted, true, true},
		"":            {statusUnknown, false, false},
		"maybe":       {statusUnknown, false, false},
		"in_use":      {statusUnknown, false, false},
		"executed, x": {statusUnknown, false, false},
	}
	for value, test := range tests {
		header := http.Header{}
		if value != "" {
			header.Set(IdempotencyKeyStatusHeader, value)
		}
		status := idempotencyKeyStatusOf(header)
		require.Equal(t, test.status, status, value)
		require.Equal(t, test.known, status.known(), value)
		require.Equal(t, test.honored, status.keyHonored(), value)
	}
}
