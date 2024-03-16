package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type User struct {
	EntityID        string `json:"entityId"`
	Username        string `json:"username"`
	Score           int    `json:"score"`
	NoOfGamesPlayed int    `json:"noOfGamesPlayed"`
}

var ctx = context.Background()
var redisClient *redis.Client

func createUser(username string) (*User, error) {
	entityID := uuid.New().String()
	user := &User{
		EntityID:        entityID,
		Username:        username,
		Score:           0,
		NoOfGamesPlayed: 0,
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	err = redisClient.HSet(ctx, "users", entityID, userJSON).Err()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var user User
		json.NewDecoder(r.Body).Decode(&user)

		createdUser, err := createUser(user.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(createdUser)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func updateUserScore(entityID string, score int) error {
	userJSON, err := redisClient.HGet(ctx, "users", entityID).Bytes()
	if err != nil {
		return err
	}

	var user User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return err
	}

	user.Score = score

	userJSON, err = json.Marshal(user)
	if err != nil {
		return err
	}

	err = redisClient.HSet(ctx, "users", entityID, userJSON).Err()
	if err != nil {
		return err
	}

	return nil
}

func handleUserScore(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Parse the entity ID from the request URL
		entityID := r.URL.Query().Get("entityId")
		if entityID == "" {
			http.Error(w, "Entity ID is required", http.StatusBadRequest)
			return
		}

		// Retrieve the user data from Redis based on the entity ID
		userJSON, err := redisClient.HGet(ctx, "users", entityID).Bytes()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Unmarshal the user data into a User struct
		var user User
		err = json.Unmarshal(userJSON, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the user's score in the response body
		responseData := struct {
			Score int `json:"score"`
		}{
			Score: user.Score,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)
	case http.MethodPost:
		var req struct {
			EntityID string `json:"entityId"`
			Score    int    `json:"score"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		err := updateUserScore(req.EntityID, req.Score)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	http.HandleFunc("/user", handleUser)
	// Register the combined handler for /user/score endpoint
	http.HandleFunc("/user/score", handleUserScore)

	// Start the HTTP server
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
