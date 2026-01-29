package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func respondWithJSON(response http.ResponseWriter, statusCode int, payload interface{}) {
	// Set the Content-Type response header to indicate JSON data
	response.Header().Set("Content-Type", "application/json")

	// Write the HTTP status code (e.g. 200, 404, 500)
	// This must be called before writing the response body
	response.WriteHeader(statusCode)

	// Encode the payload into JSON and write it to the response body
	// The payload can be any type that can be marshaled to JSON
	json.NewEncoder(response).Encode(payload)
}

func respondWithError(response http.ResponseWriter, statusCode int, message string) {
	// Respond with a JSON error message
	//map[type of key]type of value
	respondWithJSON(response, statusCode, map[string]string{"error": message})
}

func getIDFromPath(request *http.Request) (int, error) {
	// Extract the product ID from the URL path
	// Example URL path: /products/123
	// Split the path by "/" and get the second part as the ID
	// ["", "products", "123"]
	pathParts := strings.Split(request.URL.Path, "/")
	if len(pathParts) < 3 {
		return 0, errors.New("invalid URL path")
	}
	idString := pathParts[2]
	// Convert the ID from string to integer
	// strconv.Atoi converts string to int
	// error if the string is not a valid integer
	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, errors.New("invalid product ID")
	}
	return id, nil
}
