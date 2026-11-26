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

When the server answers with an error response, the error is a `*shared.SeatsioError` with the `Code` and `Message` of the first error. If the server returns an error response without an error code (for example, an empty `errors` array), the error is not a `SeatsioError`, but a generic error with the HTTP status and the response body.

A write request (`POST`, `PUT` or `DELETE`) that fails without an error response from the API (a connection failure, a timeout, or a response that is not a seats.io error) returns a `*shared.RequestError`. It wraps the original error, so `errors.Is`, `errors.As` and `os.IsTimeout` keep working.

For a write request, both error types also have these fields (see [Idempotency](#idempotency)):

- `IdempotencyKey`: the `Idempotency-Key` the request was sent with. It's empty for the errors of other requests.
- `IdempotencyKeyHonored`: whether the server honored that key, according to the `Idempotency-Key-Status` header of the response; `nil` when there was no response or the response didn't say.
- `Replayed`: `true` when the server didn't execute the request again, but replayed the response of an earlier request with the same key.

`seatsio.ClientSupport.IsOutcomeUnknown(err)` tells you whether the request may or may not have been executed: it's `true` when the response was lost after the request was sent (a timeout, a connection reset, an incomplete response body), when the server answered with a 5xx, when it answered `409` with error code `IDEMPOTENCY_KEY_IN_USE` after the client ran out of retries (an earlier request with the same key was still being processed), or when an automatic retry got any other error that wasn't a replay (for example a `429`, or a `403` from a proxy) after an earlier attempt of the same call may have been executed. It's `false` when the request never reached the server (connection refused, DNS or TLS failure), for every other 4xx, and for the errors of `GET` requests.

`seatsio.ClientSupport.IsSafeToRetryWithSameKey(err)` tells you whether you can make the same call again with the same key (see [Idempotency](#idempotency)) without risking executing it twice. It is `true` when the request never reached the server, and when its outcome is unknown, the server honors keys, the response wasn't a replay, and the first attempt started less than 5 minutes ago. `seatsio.ClientSupport.IdempotencyKeyOf(err)` returns that key. Both errors also have the methods `IsOutcomeUnknown()` and `IsSafeToRetryWithSameKey()`.

## Rate limiting - exponential backoff

This library supports [exponential backoff](https://en.wikipedia.org/wiki/Exponential_backoff).

When you send too many concurrent requests, the server returns an error `429 - Too Many Requests`. The client reacts to this by waiting for a while, and then retrying the request.
If the request still fails with an error `429`, it waits a little longer, and try again. By default,  this happens 5 times, before giving up (after approximately 15 seconds).

To change the maximum number of retries, call `SetMaxRetries` on the `SeatsioClient`:

```go
import (
    "github.com/seatsio/seatsio-go/v13"
)

client := seatsio.NewSeatsioClient(seatsio.EU, <COMPANY ADMIN KEY>, seatsio.ClientSupport.WorkspaceKey(<WORKSPACE PUBLIC KEY>))
err := client.SetMaxRetries(3)
```

`SetMaxRetries` limits all retries of a call together, whatever their reason (see [Automatic retries](#automatic-retries)). Passing in 0 disables retries completely. In that case, the client will never retry a failed request.


## Idempotency

**NOTE**: Idempotency is currently in a soft launch: it is enabled per company and not yet publicly available. Contact seats.io support to have it enabled for your company.

Every write request (`POST`, `PUT` and `DELETE`) is sent with an `Idempotency-Key` header (a v4 UUID, generated by the client). If idempotency is enabled for your company, the server executes a request with a given key only once, and replays the original response when it receives the same request again. The server keeps a key for 10 minutes.

Each response to a `POST` or `DELETE` that reached a route says what the server did with the key, in the `Idempotency-Key-Status` header (a `PUT`, or a request that was turned away before the route, such as a `401` or a `429`, has none): `executed`, `replayed`, `in-use` (an earlier request with the same key is still being processed), `reused` (the key was used for a different request) or `ignored` (idempotency isn't enabled, so the server executed the request without looking at the key).

If the connection fails after the server executed a request but before you received the response, you don't know whether the operation succeeded. The error contains the key, and tells you whether retrying with that key is safe:

```go
_, err := client.Events.Create(ctx, params)
// true when the request never reached the server, or when its outcome is unknown (a lost response, a 5xx,
// a 409 IDEMPOTENCY_KEY_IN_USE), the server honors idempotency keys for it, and the first attempt
// was less than 5 minutes ago. Never true for a response that was itself a replay.
if seatsio.ClientSupport.IsSafeToRetryWithSameKey(err) {
    key, _ := seatsio.ClientSupport.IdempotencyKeyOf(err)
    retryCtx := seatsio.ClientSupport.WithIdempotencyKey(ctx, key)
    _, err = client.Events.Create(retryCtx, params)
}
```

Retry promptly: after 10 minutes, the server has forgotten the key. The retried request must be identical to the original one, otherwise the server returns an error `422`. Every write request made with a context returned by `WithIdempotencyKey` sends that key, so use it for one operation only. You can also supply your own key: create it with `ClientSupport.NewIdempotencyKey()`, or validate one you already have with `ClientSupport.ParseIdempotencyKey(s)`. It must be a v4 UUID in its usual form (`xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx`, any case); anything else is rejected. A zero `IdempotencyKey` in the context makes the request fail before it is sent.

### Automatic retries

The client retries a write with the same key:

- when the server answers `429`, or `409` with error code `IDEMPOTENCY_KEY_IN_USE`;
- when the request never reached the server (connection refused, connect timeout, DNS or TLS failure);
- when the response was lost after the request was sent, but only if the server has confirmed that it honors idempotency keys. Idempotency is switched on or off for all write requests at once, so the client remembers the `Idempotency-Key-Status` of the last response that had one, whatever the operation. The first write of a client is never retried after a lost response, and no write is retried after one while the last status was `ignored`.

All these retries count against `SetMaxRetries` (default 5; 0 turns them off). Once an attempt of a request may have been executed (its response was lost, or the key was in use), the client doesn't start a retry more than 5 minutes after the first attempt, well within the lifetime of the key. When the server can't be reached at all (connection refused, connect timeout, DNS or TLS failure), the client stops retrying 30 seconds after the first attempt, so a write against an unreachable server takes at most about 40 seconds with the default settings. These limits can only stop retries earlier, never add any. A `422` and an invalid key are never retried. Each `SeatsioClient` has its own memory.

Because of that memory, create one `SeatsioClient` and share it across your application (it's safe for concurrent use by multiple goroutines), instead of creating a new one per request. A client that is created per request starts with an empty memory, so it never retries after a lost response, and idempotency can't protect you from executing an operation twice.

There is a small remaining risk: if seats.io switches idempotency off (globally, or for your company) between the last response that said keys are honored and a lost response, the retry can execute the operation a second time. The next response with `ignored` corrects the client's memory, so this window is short.

`GET` requests are retried after a `429`, and when the connection fails or times out before the complete response was received.

When the context of a request is cancelled or expires, the client stops retrying it, also when it is waiting before a retry. This applies to all requests, `GET` included.
