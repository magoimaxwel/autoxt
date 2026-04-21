package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const brevoAPI = "https://api.brevo.com/v3/contacts"

type BrevoContact struct {
	Email         string `json:"email"`
	ListIDs       []int  `json:"listIds"`
	UpdateEnabled bool   `json:"updateEnabled"`
}

func addContactToBrevo(email string) error {
	contact := BrevoContact{
		Email:         email,
		ListIDs:       []int{2},
		UpdateEnabled: true,
	}

	body, err := json.Marshal(contact)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", brevoAPI, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", getEnv("BREVO_API_KEY", ""))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Log the response
	log.Printf("Brevo response status: %d for email: %s", resp.StatusCode, email)

	if resp.StatusCode != 201 && resp.StatusCode != 204 {
		return fmt.Errorf("brevo returned status %d", resp.StatusCode)
	}

	return nil
}
