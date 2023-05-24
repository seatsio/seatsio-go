package workspaces

type Workspace struct {
	Name      string `json:"name"`
	Key       string `json:"key"`
	SecretKey string `json:"secretKey"`
	IsTest    bool   `json:"isTest"`
	IsActive  bool   `json:"isActive"`
	IsDefault bool   `json:"isDefault"`
}
