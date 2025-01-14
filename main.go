package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GitHubEvent struct {
	ID   string   `json:"id"`
	Type string   `json:"type"`
	Repo struct { // The name of the repo is a nested struct.
		Name string `json:"name"`
	} `json:"repo"`
	CreatedAt string `json:"created_at"`
}

func main() {
	fmt.Println("GitHub Activity Tracker 1.0\n")

	var username string
	fmt.Print("Enter the GitHub username to track: ")
	fmt.Scanln(&username)

	var url string = "https://api.github.com/users/" + username + "/events"

	// HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Check status code
	if response.StatusCode != http.StatusOK {
		switch response.StatusCode {
		case http.StatusNotFound:
			fmt.Println("Error: 404 Not Found - Likely invalid username.")
		case http.StatusUnauthorized:
			fmt.Println("Error: 401 Unauthorized - Authentication is required and has failed or has not been provided.")
		case http.StatusForbidden:
			fmt.Println("Error: 403 Forbidden - You do not have permission to access the requested resource.")
		case http.StatusInternalServerError:
			fmt.Println("Error: 500 Internal Server Error - The server encountered an unexpected condition.")
		case http.StatusBadRequest:
			fmt.Println("Error: 400 Bad Request - The server could not understand the request due to invalid syntax.")
		case http.StatusConflict:
			fmt.Println("Error: 409 Conflict - The request could not be completed due to a conflict with the current state of the target resource.")
		default:
			fmt.Printf("Error: %d %s\n", response.StatusCode, response.Status)
		}
		return
	}

	var events []GitHubEvent            // As defined at the start.
	err = json.Unmarshal(body, &events) // where interface is any type.

	fmt.Println("Total events: ", len(events))
	for _, event := range events {
		fmt.Printf("\nEvent ID: %v\n", event.ID)
		fmt.Printf("Event type: %v\n", event.Type)
		fmt.Printf("Repository name: %v\n", event.Repo.Name)
		fmt.Printf("Created at: %v\n", event.CreatedAt)
	}

}
