package controllers

import (
	"context"
	"encoding/json"
	// "fmt"
	"net/http"
	"strings"
	"time"

	"github.com/RG-7/m-backend/database"
	"github.com/RG-7/m-backend/helpers"
	"github.com/RG-7/m-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	collection := database.Client.Database("ttms").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check if user already exists (by email, mobile number, employee ID, or faculty code)
	var existingUser models.User
	filter := bson.M{"$or": []bson.M{
		{"email": user.Email},
		{"mobileno": user.MobileNumber},
		{"facultyCode": user.FacultyCode},
		{"employeeId": user.EmployeeID},
	}}

	err = collection.FindOne(ctx, filter).Decode(&existingUser)
	if err == nil {
		http.Error(w, "User already exists with provided email, mobile number, faculty code or employee ID", http.StatusConflict)
		return
	}

	// Hash the password before saving it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword) // Store the hashed password

	// If user doesn't exist, create new entry
	user.ID = primitive.NewObjectID()
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// func to handel login
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "appication/json")

	// step 1: parse the login request
	var Credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// decode the requrest body into the cred struct
	err := json.NewDecoder(r.Body).Decode(&Credentials)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// step 2: find user by email
	collection := database.Client.Database("ttms").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user models.User
	err = collection.FindOne(ctx, bson.M{"email": Credentials.Email}).Decode(&user)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	//step 3:
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Credentials.Password))
	if err != nil {
		http.Error(w, "Incorrect Password!", http.StatusUnauthorized)
		return
	}

	// step 4: create jwt token
	token, err := helpers.GenerateJWT(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// step 5: repsond wih the jwt
	reponse := map[string]string{"token": token}
	json.NewEncoder(w).Encode(reponse)

}

func ValidateToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
		return
	}

	tokenString := tokenParts[1]
	//fmt.Print(tokenString)

	// Verify token
	user, err := helpers.VerifyToken(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// 🔍 DEBUG: Print decoded user info
	//fmt.Println("Decoded User:", user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"msg":  "Token is valid",
		"user": user,
	})
}
