package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	asana "performance-dashboard-backend/internal/asana"
	db "performance-dashboard-backend/internal/database"

	api "performance-dashboard-backend/internal/api"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func connectDatabase() {
	err := db.ConnectMongoDB()
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
}

func fetchAsanaTasks() {
	token := os.Getenv("ASANA_TOKEN") // safer to set as env var
	projectID := os.Getenv("ASANA_PROJECT_ID")
	tasks, err := asana.FetchTasks(token, projectID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print tasks with custom fields
	for _, task := range tasks {
		fmt.Printf("ID: %s \n Task: %s (Completed: %v, Due: %s)\n", task.Gid, task.Name, task.Completed, task.DueOn)
		if task.Assignee != nil {
			fmt.Printf("  Assignee: %s - ID: %s\n", task.Assignee.Name, task.Assignee.Gid)
		}
		for _, cf := range task.CustomFields {
			fmt.Printf("  - %s: %s\n", cf.Name, cf.DisplayValue)
		}
		fmt.Println()
	}
}

func main() {
	loadEnv()
	connectDatabase()

	// // Define startDate and endDate for the query
	// startDate := time.Date(2025, 07, 1, 0, 0, 0, 0, time.UTC)  // 7 days ago
	// endDate := time.Date(2025, 07, 2, 23, 59, 59, 0, time.UTC) // today

	// performancePoint, err := db_handler.GetPerformancePointForTeam(os.Getenv("MONGO_URI"), os.Getenv("MONGODB_NAME"), os.Getenv("MONGODB_COLLECTION_COMPLETED_TASK"), constants.Video, startDate, endDate)
	// if err != nil {
	// 	log.Fatal("Error fetching performance point:", err)
	// }
	// fmt.Printf("Performance Point: %+v\n", performancePoint)

	api.Init()
	log.Fatal(http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), nil))
}
