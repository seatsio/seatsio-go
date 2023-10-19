package seatsio

import (
	"errors"
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v6/charts"
	"github.com/seatsio/seatsio-go/v6/events"
	"github.com/seatsio/seatsio-go/v6/holdtokens"
	"github.com/seatsio/seatsio-go/v6/reports"
	"github.com/seatsio/seatsio-go/v6/seasons"
	"github.com/seatsio/seatsio-go/v6/shared"
	"github.com/seatsio/seatsio-go/v6/workspaces"
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
