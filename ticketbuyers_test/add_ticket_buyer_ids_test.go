package ticketbuyers_test

import (
	"github.com/google/uuid"
	"github.com/seatsio/seatsio-go/v12"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/seatsio/seatsio-go/v12/ticketbuyers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCanAddTicketBuyerIds(t *testing.T) {
	t.Parallel()

	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	ticketBuyerId1 := uuid.New()
	ticketBuyerId2 := uuid.New()
	ticketBuyerId3 := uuid.New()
	ticketBuyerIds := []uuid.UUID{ticketBuyerId1, ticketBuyerId2, ticketBuyerId3}

	result, err := client.TicketBuyers.Add(test_util.RequestContext(), &ticketbuyers.TicketBuyerParams{Ids: ticketBuyerIds})
	require.NoError(t, err)
	require.Contains(t, result.Added, ticketBuyerId1, ticketBuyerId2, ticketBuyerId3)
	require.Empty(t, result.AlreadyPresent)
}

func TestCanAddTicketBuyerIds_ListWithDuplicates(t *testing.T) {
	t.Parallel()

	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	ticketBuyerId1 := uuid.New()
	ticketBuyerId2 := uuid.New()
	ticketBuyerIds := []uuid.UUID{ticketBuyerId1, ticketBuyerId1, ticketBuyerId1, ticketBuyerId2, ticketBuyerId2}

	result, err := client.TicketBuyers.Add(test_util.RequestContext(), &ticketbuyers.TicketBuyerParams{Ids: ticketBuyerIds})
	require.NoError(t, err)
	require.Contains(t, result.Added, ticketBuyerId1, ticketBuyerId2)
	require.Len(t, result.Added, 2)
	require.Empty(t, result.AlreadyPresent)
}

func TestCanAddTicketBuyerIds_SameIdDoesNotGetAddedTwice(t *testing.T) {
	t.Parallel()

	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	ticketBuyerId1 := uuid.New()
	ticketBuyerId2 := uuid.New()
	ticketBuyerIds := []uuid.UUID{ticketBuyerId1, ticketBuyerId2}

	result, err := client.TicketBuyers.Add(test_util.RequestContext(), &ticketbuyers.TicketBuyerParams{Ids: ticketBuyerIds})
	require.NoError(t, err)
	require.Contains(t, result.Added, ticketBuyerId1, ticketBuyerId2)
	require.Empty(t, result.AlreadyPresent)

	result2, err2 := client.TicketBuyers.Add(test_util.RequestContext(), &ticketbuyers.TicketBuyerParams{Ids: []uuid.UUID{ticketBuyerId1}})
	require.NoError(t, err2)
	require.Contains(t, result2.AlreadyPresent, ticketBuyerId1)
	require.Len(t, result2.AlreadyPresent, 1)
}

func TestNoIdsResultsInError(t *testing.T) {
	t.Parallel()

	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	_, err := client.TicketBuyers.Add(test_util.RequestContext(), &ticketbuyers.TicketBuyerParams{Ids: []uuid.UUID{}})
	require.EqualError(t, err, "#/ids: expected minimum item count: 1, found: 0")
}
