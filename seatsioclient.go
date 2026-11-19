package seatsio

import (
	"context"
	"errors"
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v13/charts"
	"github.com/seatsio/seatsio-go/v13/eventlog"
	"github.com/seatsio/seatsio-go/v13/events"
	"github.com/seatsio/seatsio-go/v13/holdtokens"
	"github.com/seatsio/seatsio-go/v13/reports"
	"github.com/seatsio/seatsio-go/v13/seasons"
	"github.com/seatsio/seatsio-go/v13/shared"
	"github.com/seatsio/seatsio-go/v13/ticketbuyers"
	"github.com/seatsio/seatsio-go/v13/workspaces"
)

type seatsioClientNS struct{}

var ClientSupport seatsioClientNS

const baseUrlStart = "https://api-"
const baseUrlEnd = ".seatsio.net"

const (
	EU string = baseUrlStart + "eu" + baseUrlEnd
	NA string = baseUrlStart + "na" + baseUrlEnd
	SA string = baseUrlStart + "sa" + baseUrlEnd
	OC string = baseUrlStart + "oc" + baseUrlEnd
)

type SeatsioClient struct {
	baseUrl      string
	secretKey    string
	apiClient    *req.Client
	workspaceKey string
	Workspaces   *workspaces.Workspaces
	Charts       *charts.Charts
	Events       *events.Events
	HoldTokens   *holdtokens.HoldTokens
	ChartReports *reports.ChartReports
	EventReports *reports.EventReports
	UsageReports *reports.UsageReports
	Channels     *events.Channels
	Seasons      *seasons.Seasons
	EventLog     *eventlog.EventLog
	TicketBuyers *ticketbuyers.TicketBuyers
}

func NewSeatsioClient(baseUrl string, secretKey string, additionalHeaders ...shared.AdditionalHeader) *SeatsioClient {
	apiClient := shared.ApiClient(secretKey, baseUrl, additionalHeaders...)
	client := &SeatsioClient{
		baseUrl:    baseUrl,
		secretKey:  secretKey,
		apiClient:  apiClient,
		Workspaces: &workspaces.Workspaces{Client: apiClient},
		Charts: &charts.Charts{
			Client:  apiClient,
			Archive: &charts.Archive{Client: apiClient},
		},
		Events:       &events.Events{Client: apiClient},
		HoldTokens:   &holdtokens.HoldTokens{Client: apiClient},
		ChartReports: &reports.ChartReports{Client: apiClient},
		EventReports: &reports.EventReports{Client: apiClient},
		UsageReports: &reports.UsageReports{Client: apiClient},
		Channels:     &events.Channels{Client: apiClient},
		Seasons:      &seasons.Seasons{Client: apiClient},
		EventLog:     &eventlog.EventLog{Client: apiClient},
		TicketBuyers: &ticketbuyers.TicketBuyers{Client: apiClient},
	}
	return client
}

func (c *SeatsioClient) SetMaxRetries(count int) error {
	if count < 0 {
		return errors.New("retry count must not be negative")
	}
	c.apiClient.SetCommonRetryCount(count)
	return nil
}

// SetRetryWritesOnConnectionFailure controls whether POST, PUT and DELETE requests are retried, with the same
// idempotency key, when the connection fails or times out before the complete response was received. Enable it only if
// idempotency is enabled for your company. 429 and 409 IDEMPOTENCY_KEY_IN_USE responses are always retried.
func (c *SeatsioClient) SetRetryWritesOnConnectionFailure(enabled bool) {
	shared.SetRetryWritesOnConnectionFailure(c.apiClient, enabled)
}

type IdempotencyKey = shared.IdempotencyKey

func (seatsioClientNS) NewIdempotencyKey() IdempotencyKey {
	return shared.NewIdempotencyKey()
}

func (seatsioClientNS) ParseIdempotencyKey(s string) (IdempotencyKey, error) {
	return shared.ParseIdempotencyKey(s)
}

func (seatsioClientNS) WithIdempotencyKey(ctx context.Context, key IdempotencyKey) context.Context {
	return shared.WithIdempotencyKey(ctx, key)
}

func (seatsioClientNS) IdempotencyKeyOf(err error) (IdempotencyKey, bool) {
	return shared.IdempotencyKeyOf(err)
}

func (seatsioClientNS) IsOutcomeUnknown(err error) bool {
	return shared.IsOutcomeUnknown(err)
}

func (seatsioClientNS) WorkspaceKey(key string) shared.AdditionalHeader {
	return shared.WithAdditionalHeader("X-Workspace-Key", key)
}
