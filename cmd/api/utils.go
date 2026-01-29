package main

import (
	"lesson-proj/internal/handlers"
	"net/http"
)

func methodHandler(handlerFunc http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if request.Method != allowedMethod {
			http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlerFunc(response, request)
	}
}

func productIDHandler(handlers *handlers.ProductHandler) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			handlers.GetProductByID(response, request)
		case http.MethodPut:
			handlers.UpdateProduct(response, request)
		case http.MethodDelete:
			handlers.DeleteProduct(response, request)
		default:
			http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
func userIDHandler(handlers *handlers.UserHandler) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			handlers.GetUserByID(response, request)
		case http.MethodPut:
			handlers.UpdateUser(response, request)
		case http.MethodDelete:
			handlers.DeleteUser(response, request)
		default:
			http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		} 
	}
}
