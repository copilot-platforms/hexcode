package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	activityProtocol, err := NewActivityLogProtocol()
	if err != nil {
		panic(err)
	}

	fmt.Println(activityProtocol)

	assetsPath := "./web/dist"
	isLocal := os.Getenv("IS_LOCAL")
	if isLocal == "" {
		assetsPath = "/app/web/dist"
	}

	handlers := APIHandlers{
		ActivityLogProtocol: activityProtocol,
	}

	http.Handle("/", http.FileServer(http.Dir(assetsPath)))
	http.HandleFunc("/data", handlers.GetActivityStats)
	http.HandleFunc("/events", handlers.WebhookEvents)

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type APIHandlers struct {
	ActivityLogProtocol *ActivityLogProtocol
}

func (h *APIHandlers) GetActivityStats(w http.ResponseWriter, r *http.Request) {
	result := ActivityStatsResponse{
		Data: []ActivityStatsInfo{
			{
				Type:  ActivityStatsTypePie,
				Title: "Signup source",
				Data: []ActivityStatsDataPoint{
					{
						Key:   "abc",
						Label: "Direct signup",
						Count: 12,
					},
					{
						Key:   "xyz",
						Label: "Invited",
						Count: 25,
					},
				},
			},
			{
				Type:  ActivityStatsTypeLine,
				Title: "Signups in last 7 days",
				Data: []ActivityStatsDataPoint{
					{
						Key:   "mon",
						Label: "Mon",
						Count: 12,
					},
					{
						Key:   "tue",
						Label: "Tues",
						Count: 25,
					},
					{
						Key:   "wed",
						Label: "Wed",
						Count: 20,
					},
					{
						Key:   "thur",
						Label: "Thurs",
						Count: 35,
					},
					{
						Key:   "fri",
						Label: "Fri",
						Count: 41,
					},
					{
						Key:   "sat",
						Label: "Sat",
						Count: 25,
					},
					{
						Key:   "sun",
						Label: "Sun",
						Count: 43,
					},
				},
			},
			{
				Type:  ActivityStatsTypeBarSingle,
				Title: "Activity per portal",
				Data: []ActivityStatsDataPoint{
					{
						Key:   "clients-deleted",
						Label: "Clients deleted",
						Count: 12,
					},
					{
						Key:   "new-clients-activated",
						Label: "New clients activated",
						Count: 25,
					},
					{
						Key:   "forms-submitted",
						Label: "Forms submitted",
						Count: 20,
					},
					{
						Key:   "files-admin",
						Label: "Files by admin",
						Count: 35,
					},
					{
						Key:   "files-clients",
						Label: "Files by clients",
						Count: 41,
					},
					{
						Key:   "links-admin",
						Label: "Links by admin",
						Count: 35,
					},
					{
						Key:   "links-clients",
						Label: "Links by clients",
						Count: 41,
					},
					{
						Key:   "messages-admin",
						Label: "Messages by admin",
						Count: 35,
					},
					{
						Key:   "messages-clients",
						Label: "Messages by clients",
						Count: 41,
					},
				},
			},
			{
				Type:  ActivityStatsTypeBarMulti,
				Title: "Activity per client",
				Data: []ActivityStatsDataPoint{
					{
						Key:   "john-doe",
						Label: "Files",
						Count: 10,
					},
					{
						Key:   "john-doe",
						Label: "Links",
						Count: 2,
					},
					{
						Key:   "john-doe",
						Label: "Forms",
						Count: 19,
					},
					{
						Key:   "john-doe",
						Label: "Messages",
						Count: 8,
					},
					{
						Key:   "jane-doe",
						Label: "Files",
						Count: 15,
					},
					{
						Key:   "jane-doe",
						Label: "Links",
						Count: 10,
					},
					{
						Key:   "jane-doe",
						Label: "Forms",
						Count: 10,
					},
					{
						Key:   "jane-doe",
						Label: "Messages",
						Count: 8,
					},
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

type WebhookPayload struct {
	EventType string                 `json:"eventType"`
	Created   string                 `json:"created"`
	Object    string                 `json:"object"`
	Data      map[string]interface{} `json:"data"`
}

func (h *APIHandlers) WebhookEvents(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestPayload WebhookPayload
	err := decoder.Decode(&requestPayload)
	if err != nil {
		fmt.Println("not able to decode payload")
		return
	}

	activityLog := ActivityLog{
		EventType:  requestPayload.EventType,
		CreateDate: time.Now(),
	}
	switch requestPayload.EventType {
	case "client.created":
		activityLog.CreatedBy = CreateTypeClient
		inviteURL, ok := requestPayload.Data["inviteUrl"]
		if ok && strings.Contains(inviteURL.(string), "/u/") {
			activityLog.CreatedBy = CreateTypeAdmin
		}

		activityLog.UserID = parseField(requestPayload.Data, "id")
	case "client.deleted":
		activityLog.CreatedBy = CreateTypeAdmin
		activityLog.UserID = parseField(requestPayload.Data, "id")
	case "client.activated":
		activityLog.CreatedBy = CreateTypeClient
		activityLog.UserID = parseField(requestPayload.Data, "id")
	case "form_response_completed":
		activityLog.CreatedBy = CreateTypeClient
		createdBy := parseField(requestPayload.Data, "clientId")
		activityLog.UserID = createdBy
	case "file.created", "link.created", "message.sent":
		createdByFieldName := "createdBy"
		if requestPayload.EventType == "message.sent" {
			createdByFieldName = "senderId"
		}

		createdBy := parseField(requestPayload.Data, createdByFieldName)
		activityLog.UserID = createdBy
		activityLog.CreatedBy = CreateTypeAdmin
		if createdBy != "" {
			_, err := GetClient(createdBy)
			if err == nil {
				// it means that client is not found
				activityLog.CreatedBy = CreateTypeClient
			}
		}
	default:
		return
	}

	_ = h.ActivityLogProtocol.InsertActivity(activityLog)
}

func parseField(src map[string]interface{}, field string) string {
	fieldValue, ok := src[field]
	if ok {
		return fieldValue.(string)
	}

	return ""
}

func GetClient(id string) (map[string]interface{}, error) {
	apiKey := os.Getenv("API_KEY")

	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/clients/%s", os.Getenv("HOST"), id), nil)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}

	req.Header.Set("X-API-KEY", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	var result = map[string]interface{}{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	return result, err
}
