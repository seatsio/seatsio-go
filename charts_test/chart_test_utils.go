package charts

func getCategories(drawing map[string]interface{}) []interface{} {
	categories := drawing["categories"].(map[string]interface{})
	return categories["list"].([]interface{})
}
