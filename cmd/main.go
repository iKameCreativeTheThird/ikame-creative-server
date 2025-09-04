package main

import (
	"fmt"
	"log"
	"os"
	"performance-dashboard-backend/pkg/asana"

	"github.com/joho/godotenv"
)

func main() {
	// Replace with your Asana token and project ID
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	token := os.Getenv("ASANA_TOKEN") // safer to set as env var
	projectID := os.Getenv("ASANA_PROJECT_ID")

	tasks, err := asana.FetchTasks(token, projectID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print tasks with custom fields
	for _, task := range tasks {
		fmt.Printf("Task: %s (Completed: %v, Due: %s)\n", task.Name, task.Completed, task.DueOn)
		if task.Assignee != nil {
			fmt.Printf("  Assignee: %s\n", task.Assignee.Name)
		}
		for _, cf := range task.CustomFields {
			fmt.Printf("  - %s: %s\n", cf.Name, cf.DisplayValue)
		}
		fmt.Println()
	}
}
