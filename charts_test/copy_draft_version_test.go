package charts

import (
	"context"
	"github.com/seatsio/seatsio-go/v9"
	"github.com/seatsio/seatsio-go/v9/charts"
	"github.com/seatsio/seatsio-go/v9/events"
	"github.com/seatsio/seatsio-go/v9/test_util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCopyDraftVersion(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	ctx, cancel := context.WithTimeout(test_util.RequestContext(), 5*time.Second)
	defer cancel()

	_, _ = client.Events.Create(ctx, &events.CreateEventParams{ChartKey: chartKey})
	_ = client.Charts.Update(ctx, chartKey, &charts.UpdateChartParams{Name: "newname"})

	copiedChart, err := client.Charts.CopyDraftVersion(ctx, chartKey)
	require.NoError(t, err)
	require.Equal(t, "newname (copy)", copiedChart.Name)
}
