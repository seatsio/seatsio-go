package workspaces

import (
	"context"
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v9/shared"
)

type Workspaces struct {
	Client *req.Client
}

type CreateWorkspaceParams struct {
	Name   string `json:"name"`
	IsTest bool   `json:"isTest"`
}

type UpdateWorkspaceParams struct {
	Name string `json:"name"`
}

type regenerateSecretKeyResponse struct {
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

func (workspaces Workspaces) CreateTestWorkspace(context context.Context, name string) (*Workspace, error) {
	return workspaces.createWorkspace(context, name, true)
}

func (workspaces Workspaces) CreateProductionWorkspace(context context.Context, name string) (*Workspace, error) {
	return workspaces.createWorkspace(context, name, false)
}

func (workspaces Workspaces) createWorkspace(context context.Context, name string, isTest bool) (*Workspace, error) {
	var workspace Workspace
	result, err := workspaces.Client.R().
		SetContext(context).
		SetBody(CreateWorkspaceParams{Name: name, IsTest: isTest}).
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

func (workspaces Workspaces) RegenerateSecretKey(context context.Context, key string) (*string, error) {
	var response regenerateSecretKeyResponse
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

func (workspaces Workspaces) SetDefaultWorkspace(context context.Context, key string) error {
	result, err := workspaces.Client.R().
		SetContext(context).
		SetPathParam("key", key).
		Post("/workspaces/actions/set-default/{key}")
	return shared.AssertOkNoBody(result, err)
}

func (workspaces Workspaces) Update(context context.Context, key string, Name string) error {
	result, err := workspaces.Client.R().
		SetContext(context).
		SetBody(UpdateWorkspaceParams{Name}).
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
