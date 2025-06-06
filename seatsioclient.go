package seatsio

import (
	"errors"
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v11/charts"
	"github.com/seatsio/seatsio-go/v11/eventlog"
	"github.com/seatsio/seatsio-go/v11/events"
	"github.com/seatsio/seatsio-go/v11/holdtokens"
	"github.com/seatsio/seatsio-go/v11/reports"
	"github.com/seatsio/seatsio-go/v11/seasons"
	"github.com/seatsio/seatsio-go/v11/shared"
	"github.com/seatsio/seatsio-go/v11/ticketbuyers"
	"github.com/seatsio/seatsio-go/v11/workspaces"
)

type seatsioClientNS struct {
	apiClient *req.Client
}

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
	ClientSupport.apiClient = apiClient
	return client
}

func (c *SeatsioClient) SetMaxRetries(count int) error {
	if count < 0 {
		return errors.New("retry count must not be negative")
	}
	ClientSupport.apiClient.SetCommonRetryCount(count)
	return nil
}

func (seatsioClientNS) WorkspaceKey(key string) shared.AdditionalHeader {
	return shared.WithAdditionalHeader("X-Workspace-Key", key)
}
