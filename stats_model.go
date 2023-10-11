package main

import "time"

type CreatedByType string

const (
	CreateTypeClient CreatedByType = "client"
	CreateTypeAdmin  CreatedByType = "admin"
)

type ActivityLog struct {
	ID         string        `json:"id"`
	Type       string        `json:"type"`
	UserID     string        `json:"user_id"`
	CreateBy   CreatedByType `json:"create_by"`
	CreateDate *time.Time    `json:"create_date"`
}

type ActivityStatsResponse struct {
	Data []ActivityStatsInfo `json:"data"`
}

type ActivityStatsType string

const (
	ActivityStatsTypePie       ActivityStatsType = "pie"
	ActivityStatsTypeLine      ActivityStatsType = "line"
	ActivityStatsTypeBarSingle ActivityStatsType = "bar-single"
	ActivityStatsTypeBarMulti  ActivityStatsType = "bar-multiple"
)

type ActivityStatsInfo struct {
	Type  ActivityStatsType        `json:"type"`
	Title string                   `json:"title"`
	Data  []ActivityStatsDataPoint `json:"data"`
}

type ActivityStatsDataPoint struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Count int    `json:"count"`
}
