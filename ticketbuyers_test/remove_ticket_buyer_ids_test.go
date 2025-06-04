package ticketbuyers

import (
	"github.com/google/uuid"
	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/seatsio/seatsio-go/v11/ticketbuyers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCanRemoveTicketBuyerIds(t *testing.T) {
	t.Parallel()

	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	ticketBuyerId1 := uuid.New()
	ticketBuyerId2 := uuid.New()
	ticketBuyerIds := []uuid.UUID{ticketBuyerId1, ticketBuyerId2}
	ticketBuyerId3 := uuid.New()

	_, addErr := client.TicketBuyers.Add(test_util.RequestContext(), &ticketbuyers.TicketBuyerParams{Ids: ticketBuyerIds})
	require.NoError(t, addErr)

	result, err := client.TicketBuyers.Remove(test_util.RequestContext(), &ticketbuyers.TicketBuyerParams{Ids: []uuid.UUID{ticketBuyerId1, ticketBuyerId2, ticketBuyerId3}})
	require.NoError(t, err)
	require.Contains(t, result.Removed, ticketBuyerId1, ticketBuyerId2)
	require.Contains(t, result.NotPresent, ticketBuyerId3)
}

func TestNoIdsResultsInError(t *testing.T) {
	t.Parallel()

	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	_, err := client.TicketBuyers.Remove(test_util.RequestContext(), &ticketbuyers.TicketBuyerParams{Ids: []uuid.UUID{}})
	require.EqualError(t, err, "#/ids: expected minimum item count: 1, found: 0")
}
