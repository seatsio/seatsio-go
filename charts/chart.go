package charts

type Chart struct {
	Id   int64    `json:"id"`
	Key  string   `json:"key"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
