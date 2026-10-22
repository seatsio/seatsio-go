package workspaces

import (
	"context"

	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v13/shared"
)

type Workspaces struct {
	Client *req.Client
}

type CreateWorkspaceParams struct {
	Name        string      `json:"name"`
	IsTest      bool        `json:"isTest"`
	ForSaleMode ForSaleMode `json:"forSaleMode,omitempty"`
}

type UpdateWorkspaceParams struct {
	Name        string      `json:"name"`
	ForSaleMode ForSaleMode `json:"forSaleMode,omitempty"`
}

type ForSaleMode string

const (
	RATE_LIMITED       ForSaleMode = "RATE_LIMITED"
	UNLIMITED_BILLABLE ForSaleMode = "UNLIMITED_BILLABLE"
)

type WorkspaceOption func(options *workspaceOptions)

type workspaceOptions struct {
	forSaleMode ForSaleMode
}

type RemoveSecretKeyParams struct {
	SecretKey string `json:"secretKey"`
}

type manageSecretKeyResponse struct {
	SecretKey string `json:"secretKey"`
}

type WorkspaceStatus string

type workspaceNS struct{}

var WorkspaceSupport workspaceNS

const (
	All      WorkspaceStatus = ""
	Active   WorkspaceStatus = "/active"
	Inactive WorkspaceStatus = "/inactive"
)

func (workspaces Workspaces) CreateTestWorkspace(context context.Context, name string, opts ...WorkspaceOption) (*Workspace, error) {
	return workspaces.createWorkspace(context, name, true, opts)
}

func (workspaces Workspaces) CreateProductionWorkspace(context context.Context, name string, opts ...WorkspaceOption) (*Workspace, error) {
	return workspaces.createWorkspace(context, name, false, opts)
}

func (workspaces Workspaces) createWorkspace(context context.Context, name string, isTest bool, opts []WorkspaceOption) (*Workspace, error) {
	var workspace Workspace
	options := toWorkspaceOptions(opts)
	result, err := workspaces.Client.R().
		SetContext(context).
		SetBody(CreateWorkspaceParams{Name: name, IsTest: isTest, ForSaleMode: options.forSaleMode}).
		SetSuccessResult(&workspace).
		Post("/workspaces")
	return shared.AssertOk(result, err, &workspace)
}

func (workspaces Workspaces) Activate(context context.Context, key string) error {
	result, err := workspaces.Client.R().
		SetContext(context).
		SetPathParam("key", key).
		Post("/workspaces/{key}/actions/activate")
	return shared.AssertOkNoBody(result, err)
}

func (workspaces Workspaces) Deactivate(context context.Context, key string) error {
	result, err := workspaces.Client.R().
		SetContext(context).
		SetPathParam("key", key).
		Post("/workspaces/{key}/actions/deactivate")
	return shared.AssertOkNoBody(result, err)
}

func (workspaces Workspaces) Delete(context context.Context, key string) error {
	result, err := workspaces.Client.R().
		SetContext(context).
		SetPathParam("key", key).
		Delete("/workspaces/{key}")
	return shared.AssertOkNoBody(result, err)
}

// Deprecated: Use AddSecretKey and RemoveSecretKey instead
func (workspaces Workspaces) RegenerateSecretKey(context context.Context, key string) (*string, error) {
	var response manageSecretKeyResponse
	result, err := workspaces.Client.R().
		SetContext(context).
		SetPathParam("key", key).
		SetSuccessResult(&response).
		Post("/workspaces/{key}/actions/regenerate-secret-key")
	ok, err := shared.AssertOk(result, err, &response)
	if err == nil {
		return &ok.SecretKey, nil
	} else {
		return nil, err
	}
}

func (workspaces Workspaces) AddSecretKey(context context.Context, key string) (*string, error) {
	var response manageSecretKeyResponse
	result, err := workspaces.Client.R().
		SetContext(context).
		SetPathParam("key", key).
		SetSuccessResult(&response).
		Post("/workspaces/{key}/actions/add-secret-key")
	ok, err := shared.AssertOk(result, err, &response)
	if err == nil {
		return &ok.SecretKey, nil
	} else {
		return nil, err
	}
}

func (workspaces Workspaces) RemoveSecretKey(context context.Context, key string, secretKeyToRemove string) error {
	var response manageSecretKeyResponse
	result, err := workspaces.Client.R().
		SetContext(context).
		SetPathParam("key", key).
		SetBody(RemoveSecretKeyParams{secretKeyToRemove}).
		SetSuccessResult(&response).
		Post("/workspaces/{key}/actions/remove-secret-key")
	return shared.AssertOkNoBody(result, err)
}

func (workspaces Workspaces) SetDefaultWorkspace(context context.Context, key string) error {
	result, err := workspaces.Client.R().
		SetContext(context).
		SetPathParam("key", key).
		Post("/workspaces/actions/set-default/{key}")
	return shared.AssertOkNoBody(result, err)
}

func (workspaces Workspaces) Update(context context.Context, key string, Name string, opts ...WorkspaceOption) error {
	options := toWorkspaceOptions(opts)
	result, err := workspaces.Client.R().
		SetContext(context).
		SetBody(UpdateWorkspaceParams{Name: Name, ForSaleMode: options.forSaleMode}).
		SetPathParam("key", key).
		Post("/workspaces/{key}")
	return shared.AssertOkNoBody(result, err)
}

func (workspaces Workspaces) Retrieve(context context.Context, key string) (*Workspace, error) {
	var workspace Workspace
	result, err := workspaces.Client.R().
		SetContext(context).
		SetPathParam("key", key).
		SetSuccessResult(&workspace).
		Get("/workspaces/{key}")
	return shared.AssertOk(result, err, &workspace)
}

func (workspaces Workspaces) ListAll(context context.Context, status WorkspaceStatus, opts ...shared.PaginationParamsOption) ([]Workspace, error) {
	return workspaces.lister(context, status).All(opts...)
}

func (workspaces Workspaces) ListFirstPage(context context.Context, status WorkspaceStatus, opts ...shared.PaginationParamsOption) (*shared.Page[Workspace], error) {
	return workspaces.lister(context, status).ListFirstPage(opts...)
}

func (workspaces Workspaces) ListPageAfter(context context.Context, id int64, status WorkspaceStatus, opts ...shared.PaginationParamsOption) (*shared.Page[Workspace], error) {
	return workspaces.lister(context, status).ListPageAfter(id, opts...)
}

func (workspaces Workspaces) ListPageBefore(context context.Context, id int64, status WorkspaceStatus, opts ...shared.PaginationParamsOption) (*shared.Page[Workspace], error) {
	return workspaces.lister(context, status).ListPageBefore(id, opts...)
}

func (workspaces Workspaces) lister(context context.Context, status WorkspaceStatus) *shared.Lister[Workspace] {
	pageFetcher := shared.PageFetcher[Workspace]{
		Client:    workspaces.Client,
		Url:       string("/workspaces" + status), // avoids URL encoding of forward slash
		UrlParams: map[string]string{},
		Context:   &context,
	}
	return &shared.Lister[Workspace]{PageFetcher: &pageFetcher}
}

func (workspaceNS) WithFilter(filterValue string) shared.PaginationParamsOption {
	return func(Params *shared.PaginationParams) {
		Params.QueryParams["filter"] = filterValue
	}
}

func (workspaceNS) WithForSaleMode(forSaleMode ForSaleMode) WorkspaceOption {
	return func(options *workspaceOptions) {
		options.forSaleMode = forSaleMode
	}
}

func toWorkspaceOptions(opts []WorkspaceOption) workspaceOptions {
	options := workspaceOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}
