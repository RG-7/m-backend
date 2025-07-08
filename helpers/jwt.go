package helpers

import (
	"context"
	"errors"
	// "fmt"
	"log"
	"os"
	"time"

	"github.com/RG-7/m-backend/database"
	"github.com/RG-7/m-backend/models"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var SECRET_KEY []byte

func SetSecret(secret string) {
	SECRET_KEY = []byte(secret)
}

// Initialize the environment variables
func InitConfig() {
	// Try loading .env file (only useful for local dev)
	if err := godotenv.Load(); err != nil {
		log.Println("✅ No .env file found (expected in production)")
	}

	// Get SECRET_KEY from environment
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		log.Fatal("❌ SECRET_KEY not found in environment variables!")
	}
	SECRET_KEY = []byte(secret)
}

func GenerateJWT(user models.User) (string, error) {
	// Define token expiration time
	expirationTime := time.Now().Add(500 * time.Hour)

	// Create JWT claims
	claims := &jwt.RegisteredClaims{
		Subject:   user.ID.Hex(),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})

	if err != nil || !token.Valid {
		return models.User{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return models.User{}, errors.New("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return models.User{}, errors.New("invalid user ID in token")
	}

	// fmt.Println("User ID from token:", userID)  // 🔍 Debug log

	// Fetch user from database
	collection := database.Client.Database("ttms").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user models.User
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return models.User{}, errors.New("invalid user ID format")
	}

	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}
