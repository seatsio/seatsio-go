# seatsio-go, the official Seats.io Go SDK

[![Build](https://github.com/seatsio/seatsio-go/workflows/Build/badge.svg)](https://github.com/seatsio/seatsio-go/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/seatsio/seatsio-go)](https://goreportcard.com/report/github.com/seatsio/seatsio-go)
![License](https://img.shields.io/github/license/seatsio/seatsio-go)
![GitHub release (with filter)](https://img.shields.io/github/v/release/seatsio/seatsio-go?sort=semver&display_name=tag)
[![Go Reference](https://pkg.go.dev/badge/github.com/seatsio/seatsio-go.svg)](https://pkg.go.dev/github.com/seatsio/seatsio-go/v13)

This is the official Go client library for the [Seats.io V2 REST API](https://docs.seats.io/docs/api-overview).

## Installing

```
require (
    github.com/seatsio/seatsio-go/v13 v13.9.0
)
```

## Usage

### General instructions

To use this library, you'll need to create a `SeatsioClient`:

```go
import (
    "github.com/seatsio/seatsio-go/v13"
)

client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
```

You can find your _workspace secret key_ in the [settings section of the workspace](https://app.seats.io/workspace-settings). It is important that you keep your _secret key_ private and not expose it in-browser calls unless it is password protected.

The region should correspond to the region of your account:

- `seatsio.EU`: Europe
- `seatsio.NA`: North-America
- `seatsio.SA`: South-America
- `seatsio.OC`: Oceania

If you're unsure about your region, have a look at your [company settings page](https://app.seats.io/company-settings).

### Creating a chart and an event

```go
import (
	"context",
    "fmt"
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/charts"
    "github.com/seatsio/seatsio-go/v13/events"
)

func CreateChartAndEvent() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    chart, _ := client.Charts.Create(<context.Context>, &charts.CreateChartParams{Name: "aChart"})
    event, _ := client.Events.Create(<context.Context>, &events.CreateEventParams{ChartKey: chart.Key})
    fmt.Printf(`Created a chart with key %s and an event with key: %s`, chart.Key, event.Key)
}
```

### Creating multiple events

```go
import (
	"context",
    "fmt"
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/charts"
    "github.com/seatsio/seatsio-go/v13/events"
)

func CreateMultipleEvents() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    chart, _ := client.Charts.Create(<context.Context>, &charts.CreateChartParams{Name: "aChart"})
	result, err := client.Events.CreateMultiple(<context.Context>, 
		chartKey,
		events.CreateMultipleEventParams{EventParams: &events.EventParams{EventKey: "event1", Date: "2023-10-18"}},
		events.CreateMultipleEventParams{EventParams: &events.EventParams{EventKey: "event2", Date: "2023-10-20"}},
	)
    for _, event := range result.Events {
        fmt.Printf(`Created an event with key: %s`, event.Key)
    }
}
```

### Booking objects

Booking an object changes its status to `booked`. Booked seats are not selectable on a rendered chart.

[https://docs.seats.io/docs/api-book-objects](https://docs.seats.io/docs/api-book-objects).

```go
import (
	"context",
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/events"
)

func BookObject() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    result, _ := client.Events.Book(<context.Context>, <AN EVENT KEY>, "A-1", "A-2")
}
```

### Booking objects that are on `HOLD`

```go
import (
	"context",
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/events"
)

func BookHeldObject() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    result, _ := client.Events.BookWithHoldToken(<context.Context>, <AN EVENT KEY>, []string{"A-1", "A-2"}, <A HOLD TOKEN>)
}
```

### Booking general admission (GA) areas

Either

```go
import (
	"context",
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/events"
)

func BookGA() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    result, _ := client.Events.Book(<context.Context>, <AN EVENT KEY>, "GA1", "GA1", "GA1")
}
```

Or:

```go
import (
	"context",
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/events"
)

func BookGA() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    result, _ := client.Events.BookWithObjectProperties(<context.Context>, event.Key, events.ObjectProperties{ObjectId: "GA1", Quantity: 3})
}
```

### Releasing objects

Releasing objects changes its status to `free`. Free seats are selectable on a rendered chart.

[https://docs.seats.io/docs/api-release-objects](https://docs.seats.io/docs/api-release-objects).

```go
import (
	"context",
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/events"
)

func ReleaseObjects() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    result, _ := client.Events.Release(<context.Context>, event.Key, "A-1", "A-2")
}
```

### Changing object status

Changes the object status to a custom status of your choice. If you need more statuses than just booked and free, you can use this to change the status of a seat, table or booth to your own custom status.

[https://docs.seats.io/docs/api-custom-object-status](https://docs.seats.io/docs/api-custom-object-status)

```go
import (
	"context",
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/events"
)

func ChangeObjectStatus() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    objects, err := client.Events.ChangeObjectStatus(<context.Context>, &events.StatusChangeParams{
        Events: []string{event.Key},
        StatusChanges: events.StatusChanges{
            Status:  "unavailable",
            Objects: []events.ObjectProperties{{ObjectId: "A-1"}, {ObjectId: "A-2"}},
        },
    })
}
```

### Listing status changes

`StatusChanges()` function returns an `events.Lister`. You can use `StatusChanges().All()` to iterate over all status changes.

```go
import (
	"context",
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/events"
)

func ListStatusChanges() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    statusChanges, err := client.Events.StatusChanges(<context.Context>, event.Key, "", "objectLabel", "desc").All(shared.Pagination.PageSize(2))
    for index, change := range statusChanges {
        //Do something with the status change
    }
}
```

You can alternatively use the paginated functions to retrieve status changes. To list status changes that comes after or before a given status change, you can use `StatusChanges().ListPageAfter()` and `StatusChanges().ListPageBefore()` functions.

```go
import (
	"context",
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/events"
)

func ListStatusChanges() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    client.Events.StatusChanges(<context.Context>, <AN EVENT KEY>).ListFirstPage(<OPTIONAL parameters>)
    client.Events.statusChanges(<context.Context>, <AN EVENT KEY>).ListPageAfter(<A STATUS CHANGE ID>)
    client.Events.statusChanges(<context.Context>, <AN EVENT KEY>).ListPageBefore(<A STATUS CHANGE ID>)
}
```  

You can also pass an optional parameter to _filter_, _sort_ or _order_ status changes. For this parameter, you can you use the helper functions of events.EventSupport.

```go
import (
	"context",
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/events"
)

func ListStatusChanges() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    support := events.EventSupport
    client.Events.StatusChanges(<context.Context>, <AN EVENT KEY>, support.WithFilter("A"), support.WithSortAsc("objectLabel", "asc")).ListFirstPage(<OPTIONAL parameters>)
    client.Events.statusChanges(<context.Context>, <AN EVENT KEY>, support.WithFilter("A"), support.WithSortDesc("objectLabel", "asc")).ListPageAfter(<A STATUS CHANGE ID>)
    client.Events.statusChanges(<context.Context>, <AN EVENT KEY>, support.WithFilter("A"), support.WithSortAsc("objectLabel", "asc")).ListPageBefore(<A STATUS CHANGE ID>)
}
```  

A combination of filter, sorting order and sorting option is also possible.

### Retrieving object category and status (and other information)

```go
import (
	"context",
    "fmt"
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/events"
)

func RetrieveObjectInformation() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    info, _ := client.Events.RetrieveObjectInfo(<<context.Context>, AN EVENT KEY>, "A-1", "A-2")

    fmt.Println(info["A-1"].CategoryKey)
    fmt.Println(info["A-1"].Label)
    fmt.Println(info["A-1"].Status)

    fmt.Println(info["A-2"].CategoryKey)
    fmt.Println(info["A-2"].Label)
    fmt.Println(info["A-2"].Status)
}
```

### Event reports

Want to know which seats of an event are booked, and which ones are free? That’s where reporting comes in handy.

The report types you can choose from are:
- status
- category label
- category key
- label
- order ID

[https://docs.seats.io/docs/api-event-reports](https://docs.seats.io/docs/api-event-reports)

```go
import (
	"context",
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/reports"
)

func GetSummary() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    report, _ := client.EventReports.SummaryByStatus(<context.Context>, <AN EVENT KEY>)
}

func GetDeepReport() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    report, _ := client.EventReports.DeepSummaryByStatus(<context.Context>, <AN EVENT KEY>)
}
```

### Listing all charts

You can list all charts using `ListAll()` function which returns an array of charts.

```go
import (
	"context",
	"fmt"
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/charts"
)

func GetAllCharts() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    retrievedCharts, _ := client.Charts.ListAll(<context.Context>)
	
	fmt.Println(retrievedCharts[0].Key)
	fmt.Println(retrievedCharts[1].Key)
	fmt.Println(retrievedCharts[2].Key)
}
```

### Listing charts page by page

E.g. to show charts in a paginated list on a dashboard.

Each page contains an `Items` array of `Chart` instances, and `NextPageStartsAfter` and `PreviousPageEndsBefore` properties. Those properties are the chart IDs after which the next page starts or the previous page ends.

```go
// ... user initially opens the screen ...

import (
	"context",
	"fmt"
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/charts"
)

func GetFirstPage() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    chartsPage, _ := client.Charts.ListFirstPage(<context.Context>)

    for _, chart := range chartsPage.Items {
        fmt.Println(chart.Key)
    }
}
```

```go
// ... user clicks on 'next page' button ...

import (
	"context",
	"fmt"
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/charts"
)

func GetNextPage() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    chartsPage, err := client.Charts.List(<context.Context>).ListPageAfter(<NextPageStartsAfter>)

    for _, chart := range chartsPage.Items {
        fmt.Println(chart.Key)
    }
}
```

```go
// ... user clicks on 'previous page' button ...

import (
	"context",
	"fmt"
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/charts"
)

func GetPreviousPage() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    chartsPage, err := client.Charts.List(<context.Context>).ListPageAfter(<PreviousPageEndsBefore>)

    for _, chart := range chartsPage.Items {
        fmt.Println(chart.Key)
    }
}
```

### Creating a workspace

```go
import (
	"context"
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/workspaces"
)

func CreateWorkspace() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <WORKSPACE SECRET KEY>)
    workspace, _ := client.Workspaces.CreateProductionWorkspace(<context.Context>, "a workspace")
}
```

### Creating a chart and an event with the company admin key

```go
import (
	"context",
	"fmt"
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/charts"
    "github.com/seatsio/seatsio-go/v13/events"
)

func UsingTheCompanyAdminKey() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <COMPANY ADMIN KEY>, seatsio.ClientSupport.WorkspaceKey(<WORKSPACE PUBLIC KEY>))
    chart, _ := client.Charts.Create(<context.Context>, &charts.CreateChartParams{Name: "aChart"})
    event, _ := client.Events.Create(<context.Context>, &events.CreateEventParams{ChartKey: chart.Key})
    fmt.Printf(`Created a chart with key %s and an event with key: %s`, chart.Key, event.Key)
}
```

### Listing categories

```go
import (
	"context",
    "fmt"
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/charts"
    "github.com/seatsio/seatsio-go/v13/events"
)

func RetrieveAndListCategories() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <COMPANY ADMIN KEY>, seatsio.ClientSupport.WorkspaceKey(<WORKSPACE PUBLIC KEY>))
    categories, err := client.Charts.ListCategories(<context.Context>, <CHART KEY>)
    for _, category := range categories {
        fmt.Println(category.Label)
    }
}
```

### Updating a category

```go
import (
	"context",
    "fmt"
    "github.com/seatsio/seatsio-go/v13"
    "github.com/seatsio/seatsio-go/v13/charts"
    "github.com/seatsio/seatsio-go/v13/events"
)

func UpdateCategory() {
    client := seatsio.NewSeatsioClient(seatsio.EU, <COMPANY ADMIN KEY>, seatsio.ClientSupport.WorkspaceKey(<WORKSPACE PUBLIC KEY>))
    err = client.Charts.UpdateCategory(<context.Context>, <CHART KEY>, <CATEGORY KEY>, charts.UpdateCategoryParams{
        Label:      "New label",
        Color:      "#bbbbbb",
        Accessible: false,
    })
}
```

## Error Handling
When an API call results in an error, the `error` returned by the function is not nil and contains the following format of information:

```json
{
  "errors": [{ "code": "RATE_LIMIT_EXCEEDED", "message": "Rate limit exceeded" }],
  "messages": ["Rate limit exceeded"],
  "requestId": "123456",
  "status": 429
}
```

If the server returns an error response without an error code (for example, an empty `errors` array), the error is not a `SeatsioError`, but a generic error with the HTTP status and the response body.

## Rate limiting - exponential backoff

This library supports [exponential backoff](https://en.wikipedia.org/wiki/Exponential_backoff).

When you send too many concurrent requests, the server returns an error `429 - Too Many Requests`. The client reacts to this by waiting for a while, and then retrying the request.
If the request still fails with an error `429`, it waits a little longer, and try again. By default,  this happens 5 times, before giving up (after approximately 15 seconds).

To change the maximum number of retries, create the `SeatsioClient` as follows:

```go
import (
    "github.com/seatsio/seatsio-go/v13"
)

client := seatsio.NewSeatsioClient(seatsio.EU, <COMPANY ADMIN KEY>, seatsio.ClientSupport.WorkspaceKey(<WORKSPACE PUBLIC KEY>)).SetMaxRetries(3)
```

Passing in 0 disables exponential backoff completely. In that case, the client will never retry a failed request.


## Idempotency

**NOTE**: Idempotency is currently in a soft launch: it is enabled per company and not yet publicly available. Contact seats.io support to have it enabled for your company.

Every `POST`, `PUT` and `DELETE` request is sent with an `Idempotency-Key` header (a v4 UUID, generated by the client). If idempotency is enabled for your company, the server executes a request with a given key only once, and replays the original response (with the header `Idempotent-Replayed: true`) when it receives the same request again. If it isn't enabled, the server ignores the key.

**Even when idempotency is enabled for your company, seats.io can temporarily switch it off for all companies (for example during an incident).** The server then ignores the key, and a retried request is executed again. The client cannot detect this, so idempotency reduces the risk of duplicate operations, but does not rule it out.

If the connection fails after the server executed a request but before you received the response, you don't know whether the operation succeeded. The error contains the key, so you can retry the operation with the same key. Only do this if idempotency is enabled for your company; otherwise the retry can execute an operation that already succeeded a second time:

```go
_, err := client.Events.Create(ctx, params)
// IsOutcomeUnknown is true for a connection failure, a timeout or a 5xx; any other error would happen again on a retry
if seatsio.ClientSupport.IsOutcomeUnknown(err) {
    if key, ok := seatsio.ClientSupport.IdempotencyKeyOf(err); ok {
        retryCtx := seatsio.ClientSupport.WithIdempotencyKey(ctx, key)
        _, err = client.Events.Create(retryCtx, params)
    }
}
```

A write that fails without an error response from the API (a connection failure, a timeout, or a response that is not a seats.io error) returns a `*shared.RequestError`. It wraps the original error, so `errors.Is`, `errors.As` and `os.IsTimeout` keep working.

The retried request must be identical to the original one, otherwise the server returns an error `422`. Use a context with an idempotency key for one operation only. You can also supply your own key: create it with `ClientSupport.NewIdempotencyKey()`, or validate one you already have with `ClientSupport.ParseIdempotencyKey(s)` (it must be a v4 UUID), and pass it to `WithIdempotencyKey`.

### Automatic retries

If the server answers `409` because an earlier request with the same key is still being processed, the client waits and retries it with the same key, like it does for `429`.

By default, the client only retries `GET` requests when the connection fails or times out before the complete response was received. To also retry `POST`, `PUT` and `DELETE` requests in that case (with the same key), enable it **only if idempotency is enabled for your company**. Otherwise, a retry could execute an operation that already succeeded a second time, and even with idempotency enabled this can happen while it is temporarily switched off (see above). Responses `429` and `409` are retried regardless of this setting:

```go
client.SetRetryWritesOnConnectionFailure(true)
```

When the context of a request is cancelled or expires, the client stops retrying it, also when it is waiting before a retry. This applies to all requests, `GET` included.
