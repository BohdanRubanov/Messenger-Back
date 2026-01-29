package main

import (
	"log"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		log.Printf("%s %s %s", request.Method, request.RemoteAddr, request.URL.Path)
		next.ServeHTTP(response, request)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "*")
		response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if request.Method == "OPTIONS" {
			response.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(response, request)
	})
}
