package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/samber/lo"
)

var API_KEY = os.Getenv("TODOIST_API_KEY")

func main() {
	if API_KEY == "" {
		fmt.Println("TODOIST_API_KEY environment variable not set. You can find your API key at https://todoist.com/app/settings/integrations/developer")
		return
	}

	// Get all tasks
	tasks := getTasks()

	if len(tasks) == 0 {
		fmt.Println("No tasks to update")
		return
	}

	fmt.Printf("Found %d tasks to update:\n\n", len(tasks))
	for _, task := range tasks {
		fmt.Println(task.Content)
	}

	fmt.Println("\nProceed? (y/n)")
	var proceed string
	fmt.Scanln(&proceed)
	if proceed != "y" {
		fmt.Println("Aborting")
		return
	}

	updateTasks(tasks)

}

type Task struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	cleaned string
}

func updateTasks(tasks []*Task) {
	for _, task := range tasks {
		body := map[string]interface{}{
			"content": task.cleaned,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", "https://api.todoist.com/rest/v2/tasks/"+task.ID, strings.NewReader(string(jsonBody)))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Authorization", "Bearer "+API_KEY)
		req.Header.Add("Content-Type", "application/json")

		fmt.Println("Updated:", task.cleaned)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error executing request:", err)
			return
		}
		if resp.StatusCode != 200 {
			fmt.Println("Error response:", resp.Status)
			msg, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			fmt.Println("Error message:", string(msg))
			return
		}
	}
}

func getTasks() []*Task {
	client := &http.Client{}

	// Sadly I can't seem to filter by "[ ]" in the query so we have to read them all.
	req, err := http.NewRequest("GET", "https://api.todoist.com/rest/v2/tasks", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+API_KEY)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error executing request:", err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("Error response:", resp.Status)
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return nil
		}
		fmt.Println("Error message:", string(msg))
		return nil
	}

	var tasks []*Task

	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		fmt.Println("Error decoding response:", err)
		return nil
	}

	// Filter and clean tasks
	tasks = lo.Filter(tasks, func(task *Task, _ int) bool {
		title := strings.TrimSpace(task.Content)
		if !strings.HasPrefix(title, "[ ]") {
			return false
		}
		title = strings.TrimSpace(title[3:])
		task.cleaned = title
		return true
	})

	return tasks
}
