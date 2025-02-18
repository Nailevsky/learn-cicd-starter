package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bootdotdev/learn-cicd-starter/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("🔍 Received request to create user")

	// Parse request body
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("❌ ERROR: Couldn't decode parameters:", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	log.Println("📩 Received user data:", params.Name)

	// Generate API key
	apiKey, err := generateRandomSHA256Hash()
	if err != nil {
		log.Println("❌ ERROR: Couldn't generate API key:", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate API key")
		return
	}

	// Create user in the database
	log.Println("📝 Executing SQL query: INSERT INTO users (id, created_at, updated_at, name, api_key)")

	err = cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New().String(),
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		Name:      params.Name,
		ApiKey:    apiKey,
	})
	if err != nil {
		log.Println("❌ ERROR: Couldn't create user in the database:", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	log.Println("✅ User successfully created:", params.Name)

	// Retrieve user to confirm creation
	user, err := cfg.DB.GetUser(r.Context(), apiKey)
	if err != nil {
		log.Println("❌ ERROR: Couldn't retrieve created user:", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	// Convert to response format
	userResp, err := databaseUserToUser(user)
	if err != nil {
		log.Println("❌ ERROR: Couldn't convert user to response format:", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert user")
		return
	}

	log.Println("🚀 Sending success response for user:", userResp.Name)
	respondWithJSON(w, http.StatusCreated, userResp)
}

func generateRandomSHA256Hash() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(randomBytes)
	hashString := hex.EncodeToString(hash[:])
	return hashString, nil
}

func (cfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {

	userResp, err := databaseUserToUser(user)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert user")
		return
	}

	respondWithJSON(w, http.StatusOK, userResp)
}
