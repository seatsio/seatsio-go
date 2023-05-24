package charts

type Chart struct {
	Id   int64    `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
