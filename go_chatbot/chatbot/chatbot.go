package chatBot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// GeneratedResponse represents the expected response format from Hugging Face API
type GeneratedResponse struct {
	GeneratedText string `json:"generated_text"`
}

// Hugging Face API Key
const huggingFaceAPIKey = "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"

// ChatbotHandler handles the incoming requests to the chatbot
func ChatbotHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Question string `json:"question"`
	}

	// Parse the incoming JSON request
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding JSON: %v", err)
		return
	}

	// Call the Hugging Face API with the user’s question
	huggingFaceURL := "https://api-inference.huggingface.co/models/facebook/blenderbot-400M-distill"
	// Prepare the payload with additional options
	payload := map[string]interface{}{
		"inputs": requestData.Question,
		"options": map[string]interface{}{
			"temperature": 0.6,
			"max_length":  30,
		},
	}

	// Marshal the payload into bytes
	//convert the payload into a json format
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to create request payload", http.StatusInternalServerError)
		log.Printf("Error marshaling payload: %v", err)
		return
	}

	// Create a new HTTP request
	// prepare the payload to the hugging face url and get a response
	req, err := http.NewRequest("POST", huggingFaceURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		log.Printf("Error creating request: %v", err)
		return
	}

	// Set the headers
	// ensures that the server can properly parse the data you're sending.
	req.Header.Set("Content-Type", "application/json")
	//provides the server with proof that you have permission to access or modify the resource you're requesting.
	req.Header.Set("Authorization", "Bearer "+huggingFaceAPIKey)

	// Send the request, The http.Client is responsible for sending the HTTP request you created earlier.
	client := &http.Client{}
	//This method performs the actual network operation, contacting the Hugging Face API endpoint
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to contact Hugging Face API", http.StatusInternalServerError)
		log.Printf("Error contacting Hugging Face API: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check for successful status code
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("API call failed with status: %s", resp.Status), http.StatusInternalServerError)
		body, _ := ioutil.ReadAll(resp.Body) // Read the body to log the error
		log.Printf("Error response body: %s", body)
		return
	}

	// read the entire contents of the response body from an HTTP request
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		log.Printf("Error reading response body: %v", err)
		return
	}

	// Check if the response is in JSON format
	if !json.Valid(body) {
		http.Error(w, "Invalid JSON response from Hugging Face API", http.StatusInternalServerError)
		log.Printf("Non-JSON response: %s", body)
		return
	}

	// Unmarshal the response from Hugging Face
	var apiResponses []GeneratedResponse
	if err := json.Unmarshal(body, &apiResponses); err != nil {
		http.Error(w, "Failed to parse Hugging Face response", http.StatusInternalServerError)
		log.Printf("Error unmarshaling response: %v", err)
		return
	}

	// Prepare the response
	var answer string
	if len(apiResponses) == 0 || apiResponses[0].GeneratedText == "" {
		answer = "Sorry, I didn't understand that."
	} else {
		answer = apiResponses[0].GeneratedText // Get the generated text from the first response
	}

	// Send Hugging Face's answer back to the frontend
	response := map[string]string{"answer": answer}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
