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
	retries      *shared.Retries
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
	return newSeatsioClient(baseUrl, secretKey, shared.NewRetries(shared.SystemClock), additionalHeaders...)
}

func newSeatsioClient(baseUrl string, secretKey string, retries *shared.Retries, additionalHeaders ...shared.AdditionalHeader) *SeatsioClient {
	apiClient := shared.ApiClientWithRetries(secretKey, baseUrl, retries, additionalHeaders...)
	client := &SeatsioClient{
		baseUrl:    baseUrl,
		secretKey:  secretKey,
		apiClient:  apiClient,
		retries:    retries,
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

// SetMaxRetries sets how often the client retries a request after the first attempt (default 5). It is the limit on
// all retries of a call together, whatever their reason; 0 disables retries.
func (c *SeatsioClient) SetMaxRetries(count int) error {
	if count < 0 {
		return errors.New("retry count must not be negative")
	}
	c.retries.SetMaxRetries(c.apiClient, count)
	return nil
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

func (seatsioClientNS) IsSafeToRetryWithSameKey(err error) bool {
	return shared.IsSafeToRetryWithSameKey(err)
}

func (seatsioClientNS) WorkspaceKey(key string) shared.AdditionalHeader {
	return shared.WithAdditionalHeader("X-Workspace-Key", key)
}
