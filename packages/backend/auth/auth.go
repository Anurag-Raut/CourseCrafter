package auth

import (
	"CourseCrafter/database"
	"CourseCrafter/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofor-little/env"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var jwtSecret = []byte(env.Get("JWT_SECRET", ""))

func GetGoogleUrl(c *gin.Context) {
	GOOGLE_CLIENT_ID := env.Get("GOOGLE_CLIENT_ID", "")
	GOOGLE_CLIENT_SECRET := env.Get("GOOGLE_CLIENT_SECRET", "")
	conf := &oauth2.Config{
		ClientID:     GOOGLE_CLIENT_ID,
		ClientSecret: GOOGLE_CLIENT_SECRET,
		RedirectURL:  env.Get("FRONTEND_URL", "http://localhost:3000"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
	url := conf.AuthCodeURL("state")
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func LoginWithGoogle(c *gin.Context) {

	var userInfo struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	var domain = env.Get("DOMAIN", "localhost")

	if err := c.BindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	user := utils.User{Name: userInfo.Name, Email: userInfo.Email, Password: userInfo.Email, ProfileImage: &userInfo.Picture}

	loggedUser, err := database.GetUserByEmail(userInfo.Email)
	var ID int
	if err != nil {
		id, err := database.AddUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create user"})
			return
		}
		ID = id
	} else {
		ID = loggedUser.Id
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": userInfo.Name,
		"id":       ID,
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate JWT token"})
		return
	}

	c.SetCookie("token", tokenString, 3600*24, "/", domain, true, true)

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s created", userInfo.Name)})
}

func GenerateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		token, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token", "authError": true})
			return
		}

		userId, err := VerifyToken(token)
		fmt.Print(err)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "authError": true})
			return
		}

		exists, _ := database.UserExists(userId)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found", "authError": true})
		}

		fmt.Println("USER ID IN MIDDLEWARE", userId)
		c.Set("userId", userId)

		c.Next()
	}
}
func VerifyToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["id"], "asdfSAD")
		userID := int(claims["id"].(float64))

		return userID, nil
	} else {
		return 0, fmt.Errorf("invalid token")
	}
}
