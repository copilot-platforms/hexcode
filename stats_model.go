package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type CreatedByType string

const (
	CreateTypeClient CreatedByType = "client"
	CreateTypeAdmin  CreatedByType = "admin"
)

type ActivityLog struct {
	ID         int           `json:"id"`
	EventType  string        `json:"event_type"`
	UserID     string        `json:"user_id"`
	CreatedBy  CreatedByType `json:"created_by"`
	CreateDate time.Time     `json:"create_date"`
}

type ActivityLogProtocol struct {
	DB *sql.DB
}

const activityLogCreateTable string = `
CREATE TABLE IF NOT EXISTS activity_logs (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    event_type TEXT,
    user_id TEXT,
    create_by TEXT,
    created_date DATETIME
);`

func NewActivityLogProtocol() (*ActivityLogProtocol, error) {
	db, err := sql.Open("sqlite3", "./data/events.db")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(activityLogCreateTable); err != nil {
		return nil, err
	}
	return &ActivityLogProtocol{
		DB: db,
	}, nil
}

func (ap *ActivityLogProtocol) InsertActivity(activity ActivityLog) error {
	insertStudentSQL := `INSERT INTO activity_logs(event_type, user_id, create_by, created_date) VALUES (?, ?, ?, ?)`
	statement, err := ap.DB.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		fmt.Printf("error while preparing db activity : %+v\n", activity)
		return err
	}
	_, err = statement.Exec(activity.EventType, activity.UserID, activity.CreatedBy, activity.CreateDate)
	if err != nil {
		fmt.Printf("error while inserting db activity : %+v\n", activity)
		return err
	}

	return nil
}

func (ap *ActivityLogProtocol) GetCountForSignups(createdBy CreatedByType) (int, error) {
	row, err := ap.DB.Query("SELECT COUNT(*) AS record_count FROM activity_logs WHERE create_by = ? AND created_date >= datetime('now', '-7 days')", createdBy)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var count int
	for row.Next() { // Iterate and fetch the records from result cursor
		row.Scan(&count)
	}

	return count, nil
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
