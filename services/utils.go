package services

import (
	"database/sql"
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

func parseTime(raw map[string]interface{}, key string) sql.NullTime {
	if _, ok := raw[key]; !ok {
		return sql.NullTime{
			Valid: false,
		}
	}
	rawTime, ok := raw[key].(map[string]interface{})
	if !ok {
		return sql.NullTime{
			Valid: false,
		}
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
		return sql.NullTime{
			Valid: false,
		}
	}
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}
