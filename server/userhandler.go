// userHandler.go

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// UserHandler handles HTTP requests related to user operations
type UserHandler struct {
	userService UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// UpdateScore updates the score for a user
func (uh *UserHandler) UpdateScore(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = uh.userService.UpdateScore(user.EntityID, user.Score)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User score updated successfully")
}
