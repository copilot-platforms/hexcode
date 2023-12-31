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
	data := []ActivityStatsInfo{}
	signups, err := h.SignupData()
	if err == nil {
		data = append(data, *signups)
	}

	if signupTrend, err := h.SigupLineChart(); err == nil {
		data = append(data, *signupTrend)
	}

	activityPerPortal, err := h.ActivityPerPortal()
	if err == nil {
		data = append(data, *activityPerPortal)
	}

	if activityPerClient, err := h.PortalPerClient(); err == nil {
		data = append(data, *activityPerClient)
	}

	result := ActivityStatsResponse{
		Data: data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *APIHandlers) ActivityPerPortal() (*ActivityStatsInfo, error) {
	deletedClients, err := h.ActivityLogProtocol.GetCountForEvents("client.deleted", CreateTypeAdmin)
	if err != nil {
		fmt.Println(err)
	}

	activatedClients, err := h.ActivityLogProtocol.GetCountForEvents("client.activated", CreateTypeClient)
	if err != nil {
		fmt.Println(err)
	}

	formsSubmitted, err := h.ActivityLogProtocol.GetCountForEvents("form_response.completed", CreateTypeClient)
	if err != nil {
		fmt.Println(err)
	}

	filesByAdmin, err := h.ActivityLogProtocol.GetCountForEvents("file.created", CreateTypeAdmin)
	if err != nil {
		fmt.Println(err)
	}

	filesByClient, err := h.ActivityLogProtocol.GetCountForEvents("file.created", CreateTypeClient)
	if err != nil {
		fmt.Println(err)
	}

	linksByAdmin, err := h.ActivityLogProtocol.GetCountForEvents("link.created", CreateTypeAdmin)
	if err != nil {
		fmt.Println(err)
	}

	linksByClient, err := h.ActivityLogProtocol.GetCountForEvents("link.created", CreateTypeClient)
	if err != nil {
		fmt.Println(err)
	}

	messagesByAdmin, err := h.ActivityLogProtocol.GetCountForEvents("message.sent", CreateTypeAdmin)
	if err != nil {
		fmt.Println(err)
	}

	messagesByClient, err := h.ActivityLogProtocol.GetCountForEvents("message.sent", CreateTypeClient)
	if err != nil {
		fmt.Println(err)
	}

	return &ActivityStatsInfo{
		Type:  ActivityStatsTypeBarSingle,
		Title: "Portal activity",
		Data: []ActivityStatsDataPoint{
			{
				Key:   "clients-deleted",
				Label: "Clients deleted",
				Count: deletedClients,
			},
			{
				Key:   "new-clients-activated",
				Label: "New clients activated",
				Count: activatedClients,
			},
			{
				Key:   "forms-submitted",
				Label: "Forms submitted",
				Count: formsSubmitted,
			},
			{
				Key:   "files-admin",
				Label: "Files by admin",
				Count: filesByAdmin,
			},
			{
				Key:   "files-clients",
				Label: "Files by clients",
				Count: filesByClient,
			},
			{
				Key:   "links-admin",
				Label: "Links by admin",
				Count: linksByAdmin,
			},
			{
				Key:   "links-clients",
				Label: "Links by clients",
				Count: linksByClient,
			},
			{
				Key:   "messages-admin",
				Label: "Messages by admin",
				Count: messagesByAdmin,
			},
			{
				Key:   "messages-clients",
				Label: "Messages by clients",
				Count: messagesByClient,
			},
		},
	}, nil
}

func (h *APIHandlers) SignupData() (*ActivityStatsInfo, error) {
	clientCount, err := h.ActivityLogProtocol.GetCountForEvents("client.created", CreateTypeClient)
	if err != nil {
		fmt.Println(err)
	}

	adminCount, err := h.ActivityLogProtocol.GetCountForEvents("client.created", CreateTypeAdmin)
	if err != nil {
		fmt.Println(err)
	}

	return &ActivityStatsInfo{
		Type:  ActivityStatsTypePie,
		Title: "Signup source",
		Data: []ActivityStatsDataPoint{
			{
				Key:   "direct-signup",
				Label: "Direct",
				Count: clientCount,
			},
			{
				Key:   "invited",
				Label: "Invited",
				Count: adminCount,
			},
		},
	}, nil
}

func (h *APIHandlers) SigupLineChart() (info *ActivityStatsInfo, err error) {
	eventCounts, err := h.ActivityLogProtocol.EventCountOverTime("client.created")
	if err != nil {
		log.Default().Println("error while gathering data", err)
		return
	}

	data := []ActivityStatsDataPoint{}
	executionDay := time.Now()

	for i, count := range eventCounts {
		eventDay := executionDay.AddDate(0, 0, -i)
		data = append(data, ActivityStatsDataPoint{
			Key:   eventDay.Weekday().String(),
			Label: eventDay.Weekday().String(),
			Count: count,
		})
	}

	info = &ActivityStatsInfo{
		Type:  ActivityStatsTypeLine,
		Title: "Signups in last 7 days",
		Data:  data,
	}
	return
}

func (h *APIHandlers) PortalPerClient() (info *ActivityStatsInfo, err error) {
	data, err := h.ActivityLogProtocol.EventCountByUser()
	if err != nil {
		return
	}

	info = &ActivityStatsInfo{
		Type:  ActivityStatsTypeBarMulti,
		Title: "Activity per client",
		Data:  data,
	}
	return
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
	case "form_response.completed":
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
		value, ok := fieldValue.(string)
		if !ok {
			fmt.Println("panic on parsing field")
			return ""
		}

		return value
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
