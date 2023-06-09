package seatsio

import (
	"github.com/seatsio/seatsio-go/charts"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/holdtokens"
	"github.com/seatsio/seatsio-go/reports"
	"github.com/seatsio/seatsio-go/seasons"
	"github.com/seatsio/seatsio-go/shared"
	"github.com/seatsio/seatsio-go/workspaces"
)

type seatsioClientNS struct{}

var ClientSupport seatsioClientNS

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
	Channels     *events.Channels
	Seasons      *seasons.Seasons
}

func NewSeatsioClient(secretKey string, baseUrl string, additionalHeaders ...shared.AdditionalHeader) *SeatsioClient {
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
		Channels:     &events.Channels{Client: apiClient},
		Seasons:      &seasons.Seasons{Client: apiClient},
	}
	return client
}

func (seatsioClientNS) WorkspaceKey(key string) shared.AdditionalHeader {
	return shared.WithAdditionalHeader("X-Workspace-Key", key)
}
