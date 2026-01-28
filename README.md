# seatsio-go, the official Seats.io Go SDK

[![Build](https://github.com/seatsio/seatsio-go/workflows/Build/badge.svg)](https://github.com/seatsio/seatsio-go/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/seatsio/seatsio-go)](https://goreportcard.com/report/github.com/seatsio/seatsio-go)
![License](https://img.shields.io/github/license/seatsio/seatsio-go)
![GitHub release (with filter)](https://img.shields.io/github/v/release/seatsio/seatsio-go?sort=semver&display_name=tag)
[![Go Reference](https://pkg.go.dev/badge/github.com/seatsio/seatsio-go.svg)](https://pkg.go.dev/github.com/seatsio/seatsio-go/v12)

This is the official Go client library for the [Seats.io V2 REST API](https://docs.seats.io/docs/api-overview).

## Installing

```
require (
    github.com/seatsio/seatsio-go/v12 v12.3.0
)
```

## Usage

### General instructions

To use this library, you'll need to create a `SeatsioClient`:

```go
import (
    "github.com/seatsio/seatsio-go/v12"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/charts"
    "github.com/seatsio/seatsio-go/v12/events"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/charts"
    "github.com/seatsio/seatsio-go/v12/events"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/events"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/events"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/events"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/events"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/events"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/events"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/events"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/events"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/events"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/events"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/reports"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/charts"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/charts"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/charts"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/charts"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/workspaces"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/charts"
    "github.com/seatsio/seatsio-go/v12/events"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/charts"
    "github.com/seatsio/seatsio-go/v12/events"
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
    "github.com/seatsio/seatsio-go/v12"
    "github.com/seatsio/seatsio-go/v12/charts"
    "github.com/seatsio/seatsio-go/v12/events"
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

## Rate limiting - exponential backoff

This library supports [exponential backoff](https://en.wikipedia.org/wiki/Exponential_backoff).

When you send too many concurrent requests, the server returns an error `429 - Too Many Requests`. The client reacts to this by waiting for a while, and then retrying the request.
If the request still fails with an error `429`, it waits a little longer, and try again. By default,  this happens 5 times, before giving up (after approximately 15 seconds).

To change the maximum number of retries, create the `SeatsioClient` as follows:

```go
import (
    "github.com/seatsio/seatsio-go/v12"
)

client := seatsio.NewSeatsioClient(seatsio.EU, <COMPANY ADMIN KEY>, seatsio.ClientSupport.WorkspaceKey(<WORKSPACE PUBLIC KEY>)).SetMaxRetries(3)
```

Passing in 0 disables exponential backoff completely. In that case, the client will never retry a failed request.
