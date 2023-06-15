package seatsio

import (
	"github.com/seatsio/seatsio-go/charts"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/holdtokens"
	reports "github.com/seatsio/seatsio-go/reports/charts"
	"github.com/seatsio/seatsio-go/shared"
	"github.com/seatsio/seatsio-go/workspaces"
)

type SeatsioClient struct {
	baseUrl      string
	secretKey    string
	Workspaces   *workspaces.Workspaces
	Charts       *charts.Charts
	Events       *events.Events
	HoldTokens   *holdtokens.HoldTokens
	ChartReports *reports.ChartReports
}

func NewSeatsioClient(secretKey string, baseUrl string) *SeatsioClient {
	apiClient := shared.ApiClient(secretKey, baseUrl)
	return &SeatsioClient{
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
	}
}
