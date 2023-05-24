package workspaces

import (
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/shared"
)

type Workspaces struct {
	Client *req.Client
}

type CreateWorkspaceParams struct {
	Name   string `json:"name"`
	IsTest bool   `json:"isTest"`
}

func (workspaces Workspaces) CreateTestWorkspace(name string) (*Workspace, error) {
	return workspaces.createWorkspace(name, true)
}

func (workspaces Workspaces) CreateProductionWorkspace(name string) (*Workspace, error) {
	return workspaces.createWorkspace(name, false)
}

func (workspaces Workspaces) createWorkspace(name string, isTest bool) (*Workspace, error) {
	var workspace Workspace
	result, err := workspaces.Client.R().
		SetBody(CreateWorkspaceParams{Name: name, IsTest: isTest}).
		SetSuccessResult(&workspace).
		Post("/workspaces")
	return shared.AssertOk(result, err, &workspace)
}
