package handlers

import (
	"encoding/json"
	"lesson-proj/internal/models"
	services "lesson-proj/internal/services/auth"

	"net/http"
)

type UserHandler struct {
	service *services.UserService
}

// NewUserHandler — factory function (constructor).
// It creates a new UserHandler object.
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (handler *UserHandler) Registration(response http.ResponseWriter, request *http.Request) {
	var input models.CreateUser
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid request payload")
		return
	}

	createdUser, err := handler.service.Registration(request.Context(), input)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(response, http.StatusCreated, createdUser)
}

func (handler *UserHandler) Authorization(response http.ResponseWriter, request *http.Request) {
	var authUser models.AuthUser 
	err := json.NewDecoder(request.Body).Decode(&authUser)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// context from request used to pass deadlines, cancelation signals, and other request-scoped values
	user, err := handler.service.Authorization(request.Context(), authUser.Email, authUser.Password)
	if err != nil {
		respondWithError(response, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(response, http.StatusOK, user)
}

func (handler *UserHandler) GetAllUsers(response http.ResponseWriter, request *http.Request) {
	users, err := handler.service.GetAllUsers(request.Context())
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}
	respondWithJSON(response, http.StatusOK, users)
}

func (handler *UserHandler) GetUserByID(response http.ResponseWriter, request *http.Request) {
	// Extract the user ID from the URL path with helper function
	id, err := getIDFromPath(request)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid user ID")
		return
	}
	user, err := handler.service.GetUserByID(request.Context(), id)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}
	if user == nil {
		respondWithError(response, http.StatusNotFound, "User not found")
		return
	}
	respondWithJSON(response, http.StatusOK, user)
}

func (handler *UserHandler) UpdateUser(response http.ResponseWriter, request *http.Request) {
	id, err := getIDFromPath(request)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var userInput models.UpdateUser
	err = json.NewDecoder(request.Body).Decode(&userInput)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedUser, err := handler.service.UpdateUser(request.Context(), id, userInput)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Failed to update user")
		return
	}
	respondWithJSON(response, http.StatusOK, updatedUser)
}

func (handler *UserHandler) DeleteUser(response http.ResponseWriter, request *http.Request) {
	id, err := getIDFromPath(request)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid user ID")
		return
	}
	err = handler.service.DeleteUser(request.Context(), id)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Failed to delete user")
		return
	}
	respondWithJSON(response, http.StatusNoContent, nil)
}
