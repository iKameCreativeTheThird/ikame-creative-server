package asana

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchTasks(token string, projectID string) ([]Task, error) {
	var allTasks []Task
	url := fmt.Sprintf("https://app.asana.com/api/1.0/projects/%s/tasks?opt_fields=name,assignee.name,completed,due_on,custom_fields&limit=50", projectID)

	for {
		// Build request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)

		// Send request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Read body
		body, _ := io.ReadAll(resp.Body)

		// Parse JSON
		var asanaResp AsanaResponse
		if err := json.Unmarshal(body, &asanaResp); err != nil {
			return nil, err
		}

		// Collect tasks
		allTasks = append(allTasks, asanaResp.Data...)

		// Check pagination
		if asanaResp.NextPage == nil {
			break
		}
		url = asanaResp.NextPage.Uri
	}

	return allTasks, nil
}
