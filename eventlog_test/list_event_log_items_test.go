package eventlog

import (
	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/charts"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestListAllEventLogItems(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chart, err := client.Charts.Create(test_util.RequestContext(), &charts.CreateChartParams{})
	require.NoError(t, err)

	err = client.Charts.Update(test_util.RequestContext(), chart.Key, &charts.UpdateChartParams{Name: "a chart"})
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	retrievedEventLogItems, err := client.EventLog.ListAll(test_util.RequestContext())
	require.NoError(t, err)

	require.Equal(t, 2, len(retrievedEventLogItems))
	require.Equal(t, "chart.created", retrievedEventLogItems[0].Type)
	require.Equal(t, "chart.published", retrievedEventLogItems[1].Type)
}

func TestEventLogItemProperties(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chart, err := client.Charts.Create(test_util.RequestContext(), &charts.CreateChartParams{})
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	retrievedEventLogItems, err := client.EventLog.ListFirstPage(test_util.RequestContext())
	require.NoError(t, err)

	eventLogItem := retrievedEventLogItems.Items[0]
	require.Greater(t, eventLogItem.Id, int64(0))
	require.Equal(t, "chart.created", eventLogItem.Type)
	require.NotNil(t, eventLogItem.Timestamp)
	require.Equal(t, map[string]any{"key": chart.Key, "workspaceKey": company.Workspace.Key}, eventLogItem.Data)
}
