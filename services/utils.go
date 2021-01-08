package services

import (
	"fmt"
	"time"
)

func parseString(raw map[string]interface{}, key string) string {
	if _, ok := raw[key]; !ok {
		return ""
	}
	if raw[key] == nil {
		return ""
	}
	return raw[key].(string)
}

func parseTime(raw map[string]interface{}, key string) string {
	if _, ok := raw[key]; !ok {
		return "0001-01-01"
	}
	rawTime, ok := raw[key].(map[string]interface{})
	if !ok {
		return "0001-01-01"
	}
	year := "1"
	month := "1"
	date := "1"
	yearRaw, ok := rawTime["year"]
	if ok && yearRaw != nil {
		year = fmt.Sprintf("%d", int(yearRaw.(float64)))
	}
	monthRaw, ok := rawTime["month"]
	if ok && monthRaw != nil {
		month = fmt.Sprintf("%d", int(monthRaw.(float64)))
	}
	dateRaw, ok := rawTime["date"]
	if ok && dateRaw != nil {
		date = fmt.Sprintf("%d", int(dateRaw.(float64)))
	}
	t, err := time.Parse("1-2-2006", fmt.Sprintf("%s-%s-%s", date, month, year))
	if err != nil {
		return "0001-01-01"
	}
	return t.Format("2006-01-02")
}
