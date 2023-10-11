package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.Handle("/", http.FileServer(http.Dir("./web/dist")))
	http.HandleFunc("/data", GetActivityStats)

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func GetActivityStats(w http.ResponseWriter, r *http.Request) {
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
						Label: "John doe|Files",
						Count: 12,
					},
					{
						Key:   "john-doe",
						Label: "John doe|Links",
						Count: 12,
					},
					{
						Key:   "john-doe",
						Label: "John doe|Forms",
						Count: 12,
					},
					{
						Key:   "john-doe",
						Label: "John doe|Messages",
						Count: 12,
					},
					{
						Key:   "jane-doe",
						Label: "Jane doe|Files",
						Count: 12,
					},
					{
						Key:   "jane-doe",
						Label: "jane doe|Links",
						Count: 12,
					},
					{
						Key:   "jane-doe",
						Label: "jane doe|Forms",
						Count: 12,
					},
					{
						Key:   "jane-doe",
						Label: "jane doe|Messages",
						Count: 12,
					},
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
