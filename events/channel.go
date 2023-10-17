package events

type Channel struct {
	Key     string   `json:"key"`
	Name    string   `json:"name"`
	Color   string   `json:"color"`
	Index   int      `json:"index"`
	Objects []string `json:"objects"`
}
