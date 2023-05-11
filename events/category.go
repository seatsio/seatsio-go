package events

import "encoding/json"

type CategoryKey struct {
	Key interface{}
}

func (categoryKey *CategoryKey) UnmarshalJSON(data []byte) error {
	var categoryKeyJson interface{}
	if err := json.Unmarshal(data, &categoryKeyJson); err != nil {
		return err
	}

	if categoryKeyAsFloat, isFloat := categoryKeyJson.(float64); isFloat {
		categoryKey.Key = int(categoryKeyAsFloat)
	} else {
		categoryKeyAsString, _ := categoryKeyJson.(string)
		categoryKey.Key = categoryKeyAsString
	}
	return nil
}

type Category struct {
	Key        CategoryKey `json:"key"`
	Label      string      `json:"label"`
	Color      string      `json:"color"`
	Accessible bool        `json:"accessible"`
}
