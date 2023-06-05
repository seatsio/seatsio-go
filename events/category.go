package events

import (
	"encoding/json"
	"strconv"
)

type CategoryKey struct {
	Key interface{}
}

//goland:noinspection GoMixedReceiverTypes
func (categoryKey CategoryKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(categoryKey.Key)
}

//goland:noinspection GoMixedReceiverTypes
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

func (categoryKey CategoryKey) KeyAsString() string {
	if categoryKey.isInt() {
		return strconv.Itoa(categoryKey.Key.(int))
	} else {
		return categoryKey.Key.(string)
	}
}

func (categoryKey CategoryKey) isInt() bool {
	if _, isInt := categoryKey.Key.(int); isInt {
		return true
	} else {
		return false
	}
}

type Category struct {
	Key        CategoryKey `json:"key"`
	Label      string      `json:"label"`
	Color      string      `json:"color"`
	Accessible bool        `json:"accessible"`
}
