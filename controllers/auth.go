package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	// "fmt"
	"net/http"
	"strings"
	"time"

	"github.com/RG-7/m-backend/database"
	"github.com/RG-7/m-backend/helpers"
	"github.com/RG-7/m-backend/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
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

	fmt.Printf("%s", Credentials.Email)
	fmt.Printf("%s", Credentials.Password)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	print("Herr...... at step 1")

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

	print("Herr...... at step 2")

	//step 3:
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Credentials.Password))
	if err != nil {
		http.Error(w, "Incorrect Password!", http.StatusUnauthorized)
		return
	}

	print("Herr...... at step 3")

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

// vlaidate token
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

// delete user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	vars := mux.Vars(r)
	id := vars["id"] // URL should be /users/{id}

	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Convert ID to MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Connect to MongoDB and delete the user
	collection := database.Client.Database("ttms").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if a user was actually deleted
	if result.DeletedCount == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	print("Herr...... at step 4")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}


// GetAllUsers handles fetching all users from the database
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    collection := database.Client.Database("ttms").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        log.Println("Error fetching users:", err) // Log the error to identify the issue
        http.Error(w, "Failed to fetch users: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var users []map[string]interface{}
    for cursor.Next(ctx) {
        var user map[string]interface{}
        if err := cursor.Decode(&user); err != nil {
            log.Println("Error decoding user:", err) // Log decoding errors
            http.Error(w, "Failed to decode user", http.StatusInternalServerError)
            return
        }
        users = append(users, user)
    }

    if err := cursor.Err(); err != nil {
        log.Println("Cursor iteration error:", err) // Log cursor iteration errors
        http.Error(w, "Failed to iterate over users: "+err.Error(), http.StatusInternalServerError)
        return
    }

    log.Println("Fetched users successfully:", users) // Log the fetched users for debugging
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{"users": users})
}