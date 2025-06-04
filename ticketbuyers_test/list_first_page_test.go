package ticketbuyers

import (
	"github.com/google/uuid"
	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/seatsio/seatsio-go/v11/ticketbuyers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListFirstPage(t *testing.T) {
	t.Parallel()

	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	ticketBuyerId1 := uuid.New()
	ticketBuyerId2 := uuid.New()
	ticketBuyerId3 := uuid.New()
	ticketBuyerIds := []uuid.UUID{ticketBuyerId1, ticketBuyerId2, ticketBuyerId3}

	_, addErr := client.TicketBuyers.Add(test_util.RequestContext(), &ticketbuyers.TicketBuyerParams{Ids: ticketBuyerIds})
	require.NoError(t, addErr)

	result, err := client.TicketBuyers.ListFirstPage(test_util.RequestContext())
	require.NoError(t, err)
	require.Contains(t, result.Items, ticketBuyerId1, ticketBuyerId2, ticketBuyerId3)
}

func TestListFirstPageWithPageSize(t *testing.T) {
	t.Parallel()

	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	ticketBuyerId1 := uuid.New()
	ticketBuyerId2 := uuid.New()
	ticketBuyerId3 := uuid.New()
	ticketBuyerIds := []uuid.UUID{ticketBuyerId1, ticketBuyerId2, ticketBuyerId3}

	_, addErr := client.TicketBuyers.Add(test_util.RequestContext(), &ticketbuyers.TicketBuyerParams{Ids: ticketBuyerIds})
	require.NoError(t, addErr)

	result, err := client.TicketBuyers.ListFirstPage(test_util.RequestContext(), 2)
	require.NoError(t, err)
	require.Contains(t, result.Items, ticketBuyerId1, ticketBuyerId2)
}
