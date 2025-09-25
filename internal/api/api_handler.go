package apihandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	db "performance-dashboard-backend/internal/database"
	"time"
)

type Response struct {
	Message string `json:"message"`
}

// Root handler
func RootHandler(w http.ResponseWriter, r *http.Request) {
	res := Response{Message: "Welcome to Go API listener!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// Example GET handler
func GetHandler(w http.ResponseWriter, r *http.Request) {
	res := Response{Message: "This is a GET endpoint"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// Example POST handler
func PostHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	res := map[string]interface{}{
		"received": data,
		"status":   "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func PostHandlerPerformancePoint(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	isTeam := r.URL.Query().Get("isTeam")
	startTimeStr := body["startDate"].(string)
	endTimeStr := body["endDate"].(string)
	startTime, _ := time.Parse(time.RFC3339, startTimeStr)
	endTime, _ := time.Parse(time.RFC3339, endTimeStr)

	fmt.Printf("isTeam: %s, startTime: %s, endTime: %s\n", isTeam, startTimeStr, endTimeStr)

	if isTeam == "true" {
		team := body["team"].(string)
		res, err := db.GetPerformancePointForTeam(os.Getenv("MONGO_URI"), os.Getenv("MONGODB_NAME"), os.Getenv("MONGODB_COLLECTION_COMPLETED_TASK"), team, startTime, endTime)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	} else {
		email := body["email"].(string)
		res, err := db.GetPerformancePointForIndividual(os.Getenv("MONGO_URI"), os.Getenv("MONGODB_NAME"), os.Getenv("MONGODB_COLLECTION_COMPLETED_TASK"), email, startTime, endTime)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}

func Init() {
	http.HandleFunc("/post/performance-point", PostHandlerPerformancePoint)
}
